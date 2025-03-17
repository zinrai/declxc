# declxc

A declarative CLI tool for managing LXC containers using YAML configuration.

## Overview

`declxc` allows you to define multiple LXC containers in a YAML file and perform bulk operations such as create, start, stop, and destroy. It management of LXC containers by using a declarative approach instead of running multiple commands manually.

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

### Network Settings

Use [zinrai/netshed](https://github.com/zinrai/netshed) , etc. to create bridge interfaces.

| Option       | Description                     | Required |
|--------------|---------------------------------|----------|
| type         | Network type (usually 'veth')   | Yes      |
| interface    | Host interface to connect to    | Yes      |
| ipv4_address | IPv4 address with CIDR notation | No       |
| ipv4_gateway | IPv4 gateway address            | No       |

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
