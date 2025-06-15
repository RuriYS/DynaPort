package main

import (
	"fmt"
	"net"

	"github.com/RuriYS/DynaPort/types"
	"github.com/RuriYS/DynaPort/utils"
	"golang.org/x/exp/slog"
)

const (
	listenPort = 10000
	bufferSize = 6 // this is enough lol
)

func main() {
	slog.Info("DynaPort is alive!")

	addr := net.UDPAddr{Port: listenPort, IP: net.ParseIP("0.0.0.0")}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		slog.Error("%s", err.Error())
	}
	slog.Info(fmt.Sprintf("Server started at %v", &addr))

	defer conn.Close()

	buffer := make([]byte, bufferSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if n < 2 {
			slog.Warn(fmt.Sprintf("Received malformed packet from %s", remoteAddr))
			continue
		}
		if err != nil {
			slog.Error("%s", err.Error())
		}

		protocol := types.TCP
		if string(buffer[:1]) == "u" {
			protocol = types.UDP
		}

		port, err := utils.ParsePort(buffer[1:n])
		if err != nil {
			slog.Error(fmt.Sprintf("Parsing port failed: %s", err.Error()))
			continue
		}

		slog.Info(fmt.Sprintf("Received %s %d from %s", protocol, port, remoteAddr.IP.To16()))
		err = utils.ForwardPort(remoteAddr.IP.To16().String(), uint16(port), protocol)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to forward port: %s", err.Error()))
		}
	}
}
