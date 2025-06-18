# RePorter
A simple, fast, efficient dynamic port forwarder for Linux proxies written in Go.

## Prerequisites
- Go 1.23+
- A Linux server with a public IP
- Private network with your peers
- Firewall (firewalld)

## Usage
```bash
RePorter --server # server mode
RePorter --client # client mode
```

## Features
- [x] Dynamic port forwarding
- [x] Server Implementations
  - [x] Scan sockets
  - [x] Parse sockets
  - [x] Manage forwarded ports
  - [x] Service daemon
  - [x] Allowlists/Denylists
  - [ ] TTL
  - [ ] Web Interface
  - [ ] CLI Interface
- [x] Client Implementations
  - [x] Scan sockets
  - [x] Broadcast open sockets
  - [ ] Broadcast closed sockets
  - [x] Service daemon
  - [ ] Web Interface
  - [ ] CLI Interface

Currently, all opened sockets are broadcasted to the proxy server, and can't be closed by the client.

## License
This project is licensed under the AGPL v3 License. See the [LICENSE](LICENSE.md) file for details.
