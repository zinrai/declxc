# declxc

A declarative CLI tool for managing LXC containers using YAML configuration.

## Overview

`declxc` allows you to define multiple LXC containers in a YAML file and perform bulk operations such as create, start, stop, and destroy. It simplifies management of LXC containers by using a declarative approach instead of running multiple commands manually.

## Features

- Declarative container configuration using YAML
- Bulk operations on multiple containers
- Network configuration support
- Debian package installation during container setup
- User account creation during container setup
- SSH public key deployment for users
- Idempotent operations (safe to run multiple times)

## Requirements

LXC tools installed on your system.

Debian GNU/Linux:

```bash
$ sudo apt install lxc
```

## Installation

```bash
$ go build -o declxc cmd/main.go
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

| Option   | Description                             | Required |
|----------|-----------------------------------------|----------|
| name     | Container name                          | Yes      |
| template | Template to use (e.g., ubuntu, debian)  | Yes      |
| release  | Release version (e.g., focal, bullseye) | Yes      |
| arch     | Architecture (e.g., amd64, arm64)       | Yes      |
| networks | Network configuration (array)           | No       |
| packages | Debian packages to install (array)      | No       |
| users    | User account configuration (array)      | No       |

### Network Settings

Use [zinrai/netshed](https://github.com/zinrai/netshed) , etc. to create bridge interfaces.

| Option       | Description                     | Required |
|--------------|---------------------------------|----------|
| type         | Network type (usually 'veth')   | Yes      |
| interface    | Host interface to connect to    | Yes      |
| ipv4_address | IPv4 address with CIDR notation | No       |
| ipv4_gateway | IPv4 gateway address            | No       |

### Package Settings

Debian packages can be automatically installed during container creation.

```yaml
packages:
  - nginx
  - vim
  - git
  - curl
```

The packages are installed using `apt-get` after the container is created but before users are configured.

### User Settings

User accounts can be automatically created during container setup.

| Option         | Description                                   | Required |
|----------------|-----------------------------------------------|----------|
| username       | Username for the account                      | Yes      |
| password       | Password for the account                      | Yes      |
| shell          | Login shell (default: /bin/bash)              | No       |
| ssh_key_files  | List of SSH public key file paths             | No       |

**SSH Key Files**: Paths are relative to the YAML file location. Store your SSH public keys in a `keys/` directory alongside your YAML configuration.

**Security Warning**: Passwords are stored in plain text in the YAML file. This tool is intended for development environments only. Do not use in production.

## How It Works

### Container Creation Process

1. **Container creation**: Uses `lxc-create` to create the container
2. **Network configuration**: Writes network settings to a separate config file
3. **Package installation**: Installs Debian packages using `apt-get`
4. **User creation**: Creates user accounts using `chroot` and system commands
5. **SSH key deployment**: Copies SSH public keys to user's `~/.ssh/authorized_keys`

### Package Installation

- Packages are installed after network configuration but before user creation
- Uses `chroot` to execute `apt-get update` and `apt-get install` in the container
- Runs with `DEBIAN_FRONTEND=noninteractive` to avoid interactive prompts
- If package installation fails, the entire container creation process stops

### User Account Creation

- Users are created after the container is created but before it's started
- Uses `chroot` to execute `useradd` and `chpasswd` in the container's filesystem
- Existing users are skipped (idempotent operation)
- Each user gets a home directory automatically created

### SSH Key Deployment

- SSH public key files are read from paths specified in the YAML configuration
- Keys are installed to `~/.ssh/authorized_keys` with proper permissions (700 for directory, 600 for file)
- File paths are resolved relative to the YAML file location
- Missing key files generate warnings but don't stop the process

## File Structure Example

```
project/
├── containers.yaml
├── keys/
│   ├── developer.pub
│   ├── admin.pub
│   └── tester.pub
└── README.md
```

## Example Configuration

See `examples/container.yaml` for a complete configuration example including network settings, package installation, user creation, and SSH key deployment.

## Important Notes

- **Development Use Only**: This tool stores passwords in plain text and is intended for development environments only
- **Root Privileges**: Requires sudo/root access to manage LXC containers
- **Idempotent Operations**: Safe to run create command multiple times - existing containers and users are skipped
- **SSH Server**: Containers must have SSH server installed and configured for SSH key authentication to work

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
