#!/usr/bin/env bash
set -euo pipefail

virsh_cmd=(virsh -c qemu:///system)
domains=$("${virsh_cmd[@]}" list --all --name)
"${virsh_cmd[@]}" pool-info default >/dev/null

for i in 1 2 3 4 5 6 7 8 9 10; do
  name="rocky$i"
  disk="$name.qcow2"

  if grep -Fxq "$name" <<< "$domains"; then
    if [[ $(LC_ALL=C "${virsh_cmd[@]}" domstate "$name") != 'shut off' ]]; then
      "${virsh_cmd[@]}" destroy "$name"
    fi
    "${virsh_cmd[@]}" undefine "$name" --nvram
  fi

  if "${virsh_cmd[@]}" vol-info --pool default "$disk" >/dev/null 2>&1; then
    "${virsh_cmd[@]}" vol-delete --pool default "$disk"
  fi
done
