# declxc

A declarative CLI tool for managing LXC containers using YAML configuration.

## Overview

`declxc` allows you to define multiple LXC containers in a YAML file and perform bulk operations such as create, start, stop, and destroy. It simplifies management of LXC containers by using a declarative approach instead of running multiple commands manually.

## Features

- Declarative container configuration using YAML
- Bulk operations on multiple containers
- Network configuration support
- User account creation during container setup
- **Sudo privileges configuration for users**
- SSH public key deployment for users
- Idempotent operations (safe to run multiple times)

## Requirements

LXC tools installed on your system.

Debian GNU/Linux:

```bash
$ sudo apt install lxc lxc-templates debootstrap
```

## Usage

Create containers:

```bash
$ sudo declxc create -f containers.yaml
```

Start containers:

```bash
$ sudo declxc start -f containers.yaml
```

Stop containers:

```bash
$ sudo declxc stop -f containers.yaml
```

Destroy containers:

```bash
$ sudo declxc destroy -f containers.yaml
```

## YAML Configuration

Create a YAML file to define your containers. Example: `examples/container.yaml`

### Container Settings

| Option          | Description                                   | Required |
|-----------------|-----------------------------------------------|----------|
| name            | Container name                                | Yes      |
| lxc_create_args | Arguments passed verbatim to `lxc-create`     | Yes      |
| networks        | Network configuration (array)                 | No       |
| users           | User account configuration (array)            | No       |

#### lxc_create_args

The value is passed through to `lxc-create` unchanged. `declxc` injects `-n <name>` from the `name` field, so do not include `-n` here:

```yaml
lxc_create_args: -t debian -- -r bookworm -a amd64
```

This becomes `lxc-create -n <name> -t debian -- -r bookworm -a amd64`. Anything `lxc-create` and its template scripts accept can be expressed here (e.g. the `download` template's `-d <dist>`).

Notes:

- The string is split on whitespace. Shell quoting and escaping are **not** interpreted, so values containing spaces are unsupported.
- `declxc` configures users by `chroot`-ing into `/var/lib/lxc/<name>/rootfs`. Arguments that change the storage location or lxcpath (e.g. `-B`, `-P`) will break that step.

### Network Settings

Use [zinrai/netshed](https://github.com/zinrai/netshed) , etc. to create bridge interfaces.

| Option       | Description                     | Required |
|--------------|---------------------------------|----------|
| type         | Network type (usually 'veth')   | Yes      |
| interface    | Host interface to connect to    | Yes      |
| ipv4_address | IPv4 address with CIDR notation | No       |
| ipv4_gateway | IPv4 gateway address            | No       |

### User Settings

User accounts can be automatically created during container setup.

| Option         | Description                                   | Required |
|----------------|-----------------------------------------------|----------|
| username       | Username for the account                      | Yes      |
| password       | Password for the account                      | Yes      |
| shell          | Login shell (default: /bin/bash)              | No       |
| sudo           | Grant sudo privileges (default: false)        | No       |
| ssh_key_files  | List of SSH public key file paths             | No       |

**Sudo Privileges**: When `sudo: true` is set, the user will be granted full sudo privileges without password requirement (NOPASSWD). The `sudo` package will be automatically installed if needed.

**SSH Key Files**: Paths are relative to the YAML file location. Store your SSH public keys in a `keys/` directory alongside your YAML configuration.

**Security Warning**: Passwords are stored in plain text in the YAML file, and users with sudo privileges have NOPASSWD access. This tool is intended for development environments only. Do not use in production.

## How It Works

### Container Creation Process

1. **Container creation**: Uses `lxc-create` to create the container
2. **Network configuration**: Writes network settings to a separate config file
3. **User creation**: Creates user accounts using `chroot` and system commands
4. **Sudo configuration**: Configures sudo privileges for specified users
5. **SSH key deployment**: Copies SSH public keys to user's `~/.ssh/authorized_keys`

### User Account Creation

- Users are created after the container is created but before it's started
- Uses `chroot` to execute `useradd` and `chpasswd` in the container's filesystem
- Existing users are skipped (idempotent operation)
- Each user gets a home directory automatically created

### Sudo Configuration

- If any user has `sudo: true`, the `sudo` package is automatically installed
- Sudo privileges are configured by creating a file in `/etc/sudoers.d/` for each user
- Users with sudo privileges can execute any command without password (NOPASSWD)
- Configuration is independent for each user

### SSH Key Deployment

- SSH public key files are read from paths specified in the YAML configuration
- Keys are installed to `~/.ssh/authorized_keys` with proper permissions (700 for directory, 600 for file)
- File paths are resolved relative to the YAML file location
- Missing key files generate warnings but don't stop the process

## Important Notes

- **Development Use Only**: This tool stores passwords in plain text and grants NOPASSWD sudo access. It is intended for development environments only
- **Root Privileges**: Requires sudo/root access to manage LXC containers
- **Idempotent Operations**: Safe to run create command multiple times - existing containers and users are skipped
- **SSH Server**: Containers must have SSH server installed and configured for SSH key authentication to work
- **Sudo Access**: Users with `sudo: true` have full administrative privileges without password requirement

## License

This project is licensed under the [MIT License](./LICENSE).
