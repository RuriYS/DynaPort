package server

import (
	"log/slog"
	"net"

	"github.com/RuriYS/RePorter/internal"
	"github.com/RuriYS/RePorter/types"
)

var config *types.Config

func Run() {
	slog.Info("[Server] starting")
	config = internal.GetConfig()
	Init()

	addr := &net.UDPAddr{Port: int(config.Server.Port), IP: net.ParseIP(config.Server.Host)}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("[Server] failed to start the server", "error", err.Error())
	}

	if conn == nil {
		return
	}

	defer conn.Close()

	slog.Info("[Server] starting listener", "addr", addr)
	StartListener(conn)
}
