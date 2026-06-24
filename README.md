# declxc

A declarative CLI tool for managing LXC containers using YAML configuration.

## Overview

`declxc` allows you to define multiple LXC containers in a YAML file and perform bulk operations such as create, start, stop, and destroy. It simplifies management of LXC containers by using a declarative approach instead of running multiple commands manually.

## Features

- Declarative container configuration using YAML
- Bulk operations on multiple containers
- Network configuration support
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

#### lxc_create_args

The value is passed through to `lxc-create` unchanged. `declxc` injects `-n <name>` from the `name` field, so do not include `-n` here:

```yaml
lxc_create_args: -t debian -- -r bookworm -a amd64
```

This becomes `lxc-create -n <name> -t debian -- -r bookworm -a amd64`. Anything `lxc-create` and its template scripts accept can be expressed here (e.g. the `download` template's `-d <dist>`).

Notes:

- The string is split on whitespace. Shell quoting and escaping are **not** interpreted, so values containing spaces are unsupported.
- `declxc` writes network configuration under `/var/lib/lxc/<name>/`. Arguments that change the lxcpath (e.g. `-P`) will break that step.

### Network Settings

Use [zinrai/netshed](https://github.com/zinrai/netshed) , etc. to create bridge interfaces.

| Option       | Description                     | Required |
|--------------|---------------------------------|----------|
| type         | Network type (usually 'veth')   | Yes      |
| interface    | Host interface to connect to    | Yes      |
| ipv4_address | IPv4 address with CIDR notation | No       |
| ipv4_gateway | IPv4 gateway address            | No       |

## How It Works

### Container Creation Process

1. **Container creation**: Uses `lxc-create` to create the container
2. **Network configuration**: Writes network settings to a separate config file

Provisioning inside the container (users, packages, services) is out of scope; run a configuration management tool against the created container instead.

## Important Notes

- **Root Privileges**: Requires sudo/root access to manage LXC containers
- **Idempotent Operations**: Safe to run create command multiple times - existing containers are skipped

## License

This project is licensed under the [MIT License](./LICENSE).
