package server

import (
	"errors"
	"log/slog"
	"net"

	"github.com/RuriYS/RePorter/internal/config"
)

var conn *net.UDPConn
var allowedIPs map[string]struct{}
var allowedPorts map[uint16]struct{}

func initialize() (err error) {
	config := config.GetConfig()
	allowedIPs = make(map[string]struct{}, len(config.Server.AllowedIPs))
	for _, ip := range config.Server.AllowedIPs {
		allowedIPs[ip] = struct{}{}
	}
	allowedPorts = make(map[uint16]struct{}, len(config.Server.AllowedPorts))
	for _, port := range config.Server.AllowedPorts {
		allowedPorts[port] = struct{}{}
	}

	addr := &net.UDPAddr{Port: int(config.Server.Port), IP: net.ParseIP(config.Server.Host)}
	conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	if conn == nil {
		return errors.New("failed to create UDP connection")
	}

	return nil
}

func Run() {
	slog.Info("[Server] initializing")
	err := initialize()
	if err != nil {
		slog.Error("[Server] failed to initialize", "error", err)
	}

	defer conn.Close()

	slog.Info("[Server] starting")
	StartListener(conn)
}
