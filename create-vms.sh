#!/usr/bin/env bash
set -euo pipefail
cd /tmp/goautomate-rockylab
for i in 1 2; do
  virsh -c qemu:///system vol-upload --pool default "rocky$i.qcow2" base.qcow2
  virt-install --connect qemu:///system --name "rocky$i" \
    --memory 2048 --vcpus 2 --cpu host-passthrough \
    --os-variant rocky9 --import \
    --disk "vol=default/rocky$i.qcow2,format=qcow2,bus=virtio" \
    --network "network=rockylab,model=virtio,mac=52:54:00:16:00:1$i" \
    --cloud-init "user-data=$PWD/rocky$i-user.yaml,meta-data=$PWD/rocky$i-meta.yaml" \
    --graphics none --noautoconsole
 done
