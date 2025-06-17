# DynaPort
A simple, fast, efficient dynamic port forwarder for Linux proxies written in Go.

## Prerequisites
- Go 1.23+
- A Linux server with a public IP
- Private network with your peers
- firewalld (I'll try to add iptables support later)

## Roadmap
- [x] Server Implementations
  - [x] Scan sockets
  - [x] Parse sockets
  - [x] Manage forwarded ports
  - [x] Service daemon
  - [ ] TTL
  - [ ] Allowlists/Denylists
  - [ ] Web Interface
  - [ ] CLI Interface
- [x] Client Implementations
  - [x] Scan sockets
  - [x] Broadcast open sockets
  - [x] Service daemon
  - [ ] Broadcast closed sockets
  - [ ] Web Interface
  - [ ] CLI Interface
- [x] Dynamic Configuration

## Usage
```bash
DynaPort --server # server mode
DynaPort --client # client mode
```

For now, it forwards every open port in the client, and no way of reverting it, but it's not permanent.

## License
This project is licensed under the AGPL v3 License. See the [LICENSE](LICENSE.md) file for details.
