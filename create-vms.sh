#!/usr/bin/env bash
set -euo pipefail
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
mkdir -p "$script_dir/vms"
cd "$script_dir/vms"
virsh_cmd=(virsh -c qemu:///system)
"${virsh_cmd[@]}" list --all >/dev/null
network_info=$(LC_ALL=C "${virsh_cmd[@]}" net-info rockylab)
if ! grep -Eq '^Active:[[:space:]]+yes[[:space:]]*$' <<< "$network_info"; then
  "${virsh_cmd[@]}" net-start rockylab
fi
# Reserva las diez direcciones tanto en la red activa como en su configuración.
python3 - <<'PY'
import subprocess
import xml.etree.ElementTree as ET

virsh = ['virsh', '-c', 'qemu:///system']
hosts = [dict(mac=f'52:54:00:16:00:{16 + i:02x}', name=f'rocky{i}',
              ip=f'10.0.16.{i + 9}') for i in range(1, 11)]
macs = {host['mac'] for host in hosts}
ips = {host['ip'] for host in hosts}
for scope, options in [('config', ['--inactive']), ('live', [])]:
    root = ET.fromstring(subprocess.check_output(
        virsh + ['net-dumpxml', 'rockylab'] + options, text=True))
    for host in root.findall('./ip/dhcp/host'):
        if host.get('mac', '').lower() in macs or host.get('ip') in ips:
            subprocess.run(virsh + ['net-update', 'rockylab', 'delete',
                'ip-dhcp-host', ET.tostring(host, encoding='unicode'),
                f'--{scope}'], check=True)
    for host in hosts:
        subprocess.run(virsh + ['net-update', 'rockylab', 'add-last',
            'ip-dhcp-host', ET.tostring(ET.Element('host', host), encoding='unicode'),
            f'--{scope}'], check=True)
PY
for i in 1 2 3 4 5 6 7 8 9 10; do
  name="rocky$i"
  printf -v mac '52:54:00:16:00:%02x' "$((16 + i))"
  if "${virsh_cmd[@]}" dominfo "$name" >/dev/null 2>&1; then
    state=$(LC_ALL=C "${virsh_cmd[@]}" domstate "$name")
    case "$state" in
      'shut off') "${virsh_cmd[@]}" start "$name" ;;
      running) printf '%s ya está en ejecución.\n' "$name" ;;
      *) printf 'Estado inesperado de %s: %s\n' "$name" "$state" >&2; exit 1 ;;
    esac
    continue
  fi
  if "${virsh_cmd[@]}" vol-info --pool default "$name.qcow2" >/dev/null 2>&1; then
    printf 'Ya existe el disco %s.qcow2 sin su máquina; revísalo antes de continuar.\n' "$name" >&2
    exit 1
  fi
  base_image=${BASE_IMAGE:-"$PWD/base.qcow2"}
  if [[ ! -f "$base_image" ]]; then
    printf 'Falta la imagen cloud Rocky Linux 9: %s. Usa BASE_IMAGE=/ruta/imagen.qcow2.\n' "$base_image" >&2
    exit 1
  fi
  capacity=$(qemu-img info --output=json "$base_image" | python3 -c 'import json,sys; print(json.load(sys.stdin)["virtual-size"])')
  cat > "$name-user.yaml" <<EOF
#cloud-config
hostname: $name
users:
  - default
  - name: rocky
    groups: [wheel]
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    lock_passwd: false
ssh_pwauth: true
chpasswd:
  expire: false
  users:
    - name: rocky
      password: rocky
      type: text
EOF
  printf 'instance-id: %s\nlocal-hostname: %s\n' "$name" "$name" > "$name-meta.yaml"
  "${virsh_cmd[@]}" vol-create-as --pool default --name "$name.qcow2" --capacity "$capacity" --format qcow2
  "${virsh_cmd[@]}" vol-upload --pool default "$name.qcow2" "$base_image"
  virt-install --connect qemu:///system --name "rocky$i" \
    --memory 2048 --vcpus 2 --cpu host-passthrough \
    --os-variant rocky9 --import \
    --disk "vol=default/rocky$i.qcow2,format=qcow2,bus=virtio" \
    --network "network=rockylab,model=virtio,mac=$mac" \
    --cloud-init "user-data=$PWD/rocky$i-user.yaml,meta-data=$PWD/rocky$i-meta.yaml" \
    --graphics none --noautoconsole
done
