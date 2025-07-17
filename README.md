# declxc

A declarative CLI tool for managing LXC containers using YAML configuration.

## Overview

`declxc` allows you to define multiple LXC containers in a YAML file and perform bulk operations such as create, start, stop, and destroy. It simplifies management of LXC containers by using a declarative approach instead of running multiple commands manually.

## Features

- Declarative container configuration using YAML
- Bulk operations on multiple containers
- Network configuration support
- User account creation during container setup
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
| users    | User account configuration (array)      | No       |

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

| Option   | Description                           | Required |
|----------|---------------------------------------|----------|
| username | Username for the account              | Yes      |
| password | Password for the account              | Yes      |
| shell    | Login shell (default: /bin/bash)      | No       |

**Security Warning**: Passwords are stored in plain text in the YAML file. This tool is intended for development environments only. Do not use in production.

## How It Works

### Container Creation Process

1. **Container creation**: Uses `lxc-create` to create the container
2. **Network configuration**: Writes network settings to a separate config file
3. **User creation**: Creates user accounts using `chroot` and system commands

### User Account Creation

- Users are created after the container is created but before it's started
- Uses `chroot` to execute `useradd` and `chpasswd` in the container's filesystem
- Existing users are skipped (idempotent operation)
- Each user gets a home directory automatically created

## Important Notes

- **Development Use Only**: This tool stores passwords in plain text and is intended for development environments only
- **Root Privileges**: Requires sudo/root access to manage LXC containers
- **Idempotent Operations**: Safe to run create command multiple times - existing containers and users are skipped

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
