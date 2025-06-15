package internal

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/RuriYS/DynaPort/types"
	"github.com/RuriYS/DynaPort/utils"
)

type RetCode uint8

const (
	OK RetCode = 0
	Allocated RetCode = 1
)

const (
	bufferSize = 6
)

func StartServer(host string, port uint16) {
	slog.Info("dynaport is alive!")

	addr := net.UDPAddr{Port: int(port), IP: net.ParseIP(host)}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to start the server: %s", err.Error()))
	}

	slog.Info(fmt.Sprintf("server started at %v", &addr))

	defer conn.Close()

	buffer := make([]byte, bufferSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if n < 2 {
			slog.Warn(fmt.Sprintf("received malformed packet from %s", remoteAddr))
			continue
		}
		if err != nil {
			slog.Error(fmt.Sprintf("failed to read packet: %s", err.Error()))
			continue
		}

		protocol := types.TCP
		if string(buffer[:1]) == "u" {
			protocol = types.UDP
		}

		port, err := utils.ParsePort(buffer[1:n])
		if err != nil {
			slog.Error(fmt.Sprintf("parsing port failed: %s", err.Error()))
			continue
		}

		slog.Info(fmt.Sprintf("received %s %d from %s", protocol, port, remoteAddr.IP.To16()))

		allocations, err := utils.GetAllocations()
		if err != nil {
			slog.Error(fmt.Sprintf("failed to get allocations: %s", err.Error()))
		}
		
		for _, alloc := range allocations {
			if alloc.Port == port {
				_, err = conn.WriteToUDP([]byte{byte(Allocated)}, remoteAddr)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to reply: %s", err.Error()))
				}
				continue
			}
		}

		err = utils.ForwardPort(remoteAddr.IP.To16().String(), uint16(port), protocol)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to forward port: %s", err.Error()))
			continue
		}

		_, err = conn.WriteToUDP([]byte{byte(OK)}, remoteAddr)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to reply: %s", err))
			continue
		}
	}
}
