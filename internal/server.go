package internal

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"

	"github.com/RuriYS/DynaPort/types"
	"github.com/RuriYS/DynaPort/utils"
)


const (
	packetSize = 3
)

func StartServer(host string, port uint16, verbose bool) {
	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	
	addr := net.UDPAddr{Port: int(port), IP: net.ParseIP(host)}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to start the server: %s", err.Error()))
	}
	
	slog.Info("dynaport is alive!")
	slog.Info(fmt.Sprintf("listening at %v", &addr))

	defer conn.Close()

	packet := make([]byte, packetSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(packet)
		if err != nil || n != 3 {
			slog.Warn(fmt.Sprintf("Invalid packet from %s\n", remoteAddr))
			continue
		}
		
		protocol := types.TCP
		if string(packet[:1]) == "u" {
			protocol = types.UDP
		}
		
		port := binary.BigEndian.Uint16(packet[1:])

		slog.Debug(fmt.Sprintf("received %s %d from %s", protocol, port, remoteAddr.IP.To16()))

		allocations, err := utils.GetAllocations()
		if err != nil {
			slog.Error(fmt.Sprintf("failed to get allocations: %s", err.Error()))
		}
		
		for _, alloc := range allocations {
			slog.Debug(fmt.Sprintf("checking port: %d\n", alloc.Port))
			if alloc.Port == port {
				conn.WriteToUDP([]byte{byte(types.Allocated)}, remoteAddr)
				slog.Debug(fmt.Sprintf("port %d is already allocated\n", port))
				continue
			}
		}

		err = utils.ForwardPort(remoteAddr.IP.To16().String(), uint16(port), protocol)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to forward port: %s", err.Error()))
			continue
		} else {
			slog.Info(fmt.Sprintf("port forwarded: %d/%s -> %s\n", port, protocol, remoteAddr.IP.To16()))
		}

		_, err = conn.WriteToUDP([]byte{byte(types.OK)}, remoteAddr)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to reply: %s", err))
			continue
		}
	}
}
