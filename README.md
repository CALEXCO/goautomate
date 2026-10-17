# goautomate

A Go CLI for running the same command on multiple servers in parallel over SSH. Target servers and credentials are defined in a YAML inventory.

Designed for simple queries and everyday tasks in lab or development environments.

## Features

- Reads target servers from a YAML inventory.
- Opens an SSH connection to each server and runs the command concurrently.
- Prints the server address and combined stdout and stderr output.
- Waits for all executions to finish unless an error terminates the process.

## Requirements

- Go 1.26.0 or later, as specified in `go.mod`.
- Network access to the target servers on TCP port 22.
- An SSH account with password authentication and permission to run the remote command.

## Installation

```sh
git clone https://github.com/CALEXCO/goautomate.git
cd goautomate
go mod download
go build -o goautomate .
```

Check the available options:

```sh
./goautomate --help
./goautomate exec --help
```

## Inventory

Create a YAML file, such as `hosts.local.yaml`, with the following structure. Replace the example values with those for your environment:

```yaml
target_servers:
  - ip: 192.0.2.10
    user: labuser
    password: "CHANGE_ME"
  - ip: 192.0.2.11
    user: labuser
    password: "CHANGE_ME"
```

| Field | Description |
| --- | --- |
| `target_servers` | List of servers on which to run the command. |
| `ip` | Server IP address or hostname. The application appends port `22`. |
| `user` | SSH username. |
| `password` | SSH password, stored in plain text. |

Keep inventories containing credentials outside the repository or add them to your local Git exclusions. The included `hosts.yaml` contains the development lab configuration; you do not need to use it to run the application.

## Usage

Run a command on every server in the inventory:

```sh
./goautomate exec --file hosts.local.yaml --command 'hostname'
```

Other query examples:

```sh
./goautomate exec --file hosts.local.yaml --command 'uptime'
./goautomate exec --file hosts.local.yaml --command 'df -h'
./goautomate exec --file hosts.local.yaml --command 'ls -la /var/log | head -n 10'
```

Single quotes ensure that the local shell passes the full command, including pipes, to the application for execution on the remote server.

You can also run directly from source:

```sh
go run . exec --file hosts.local.yaml --command 'uname -a'
```

| Option | Description |
| --- | --- |
| `--file` | Path to the YAML inventory. |
| `--command` | Command to run on each server. |
| `-h`, `--help` | Show help. |

Provide both `--file` and `--command`: the CLI does not currently validate them as required options. Results may appear in any order, and output from different servers may interleave because executions run concurrently.

## Current limitations

- Only password authentication is supported; SSH keys and SSH agents are not implemented.
- Connections always use port 22 and do not read `~/.ssh/config`.
- SSH host keys are not verified: the implementation uses `ssh.InsecureIgnoreHostKey()`.
- No timeouts, retries, or concurrency limits are configured.
- Inventory, connection, or remote execution errors can trigger a `panic` and terminate the entire process. There is no per-server result summary.
- When a remote command fails, its output is not displayed before the error terminates execution.
- Although intended for queries, the application does not enforce read-only commands: commands run with the SSH user's permissions.

## Project structure

```text
.
├── main.go                # Entry point
├── cmd/
│   ├── root.go            # Cobra root command
│   └── exec.go            # Options and concurrent execution
├── pkg/
│   ├── ssh.go             # SSH connections and remote execution
│   └── server/server.go   # Server model and YAML parsing
├── hosts.yaml             # Development lab inventory
└── create-vms.sh           # Auxiliary libvirt lab script
```

`create-vms.sh` starts existing `rocky1` through `rocky10` domains or creates them if missing. It requires access to `qemu:///system` and an existing persistent `rockylab` network with IPv4 DHCP on `10.0.16.0/24`. It configures DHCP reservations from `10.0.16.10` (`rocky1`) through `10.0.16.19` (`rocky10`), replacing conflicting reservations for these IPs or MAC addresses. Existing guests may need a DHCP lease renewal or restart to receive their new addresses. To create domains, install `virt-install`, `qemu-img` and Python 3, provide a Rocky Linux 9 cloud image as `vms/base.qcow2` or set `BASE_IMAGE=/absolute/path/image.qcow2`, and have a storage pool named `default`. Cloud-init files are generated in `vms/`, with SSH credentials `rocky` / `rocky` for newly created machines. Existing disks are never overwritten. `hosts.yaml` contains these ten addresses.

Run `./create-vms.sh` to create or start the lab. Run `./destroy-vms.sh` to force-stop and delete all ten domains and their `default/rockyN.qcow2` disks, including their data. The base image and network reservations are retained for reuse. These scripts are not required to use the CLI.

## Development

```sh
go build ./...
go vet ./...
go test ./...
```

The repository currently has no automated tests. Functional verification requires accessible SSH servers and a valid inventory.

## Author

Contact: [alexortegaleonato@gmail.com](mailto:alexortegaleonato@gmail.com).

## License

Licensed under the [MIT License](LICENSE).
