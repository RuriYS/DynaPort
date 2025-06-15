# DynaPort
A simple, fast, efficient dynamic port forwarder for Linux proxies written in Go.

## Prerequisites
- Go 1.23+
- A Linux server with a public IP
- Private network with your peers
- firewalld (I'll try to add iptables support later)

## Building
```bash
git clone https://github.com/RuriYS/dynaport
cd dynaport
go get
go build
```

## Running
```bash
./DynaPort
```
