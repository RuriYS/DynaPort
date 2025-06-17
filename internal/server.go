package internal

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"slices"

	"github.com/RuriYS/DynaPort/types"
	"github.com/RuriYS/DynaPort/utils"
)

const packetSize = 3

func StartServer() {
	addr := &net.UDPAddr{Port: int(config.Server.Port), IP: net.ParseIP(config.Server.Host)}
	conn := initializeServer(addr)
	if conn == nil {
		return
	}

	defer conn.Close()

	slog.Info(fmt.Sprintf("starting listener at %v", addr))
	handleListener(conn)
}

func initializeServer(addr *net.UDPAddr) *net.UDPConn {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("failed to start the server", "initializeServer", err.Error())
		return nil
	}

	slog.Info("dynaport is alive!")
	return conn
}

func handleListener(conn *net.UDPConn) {
	packet := make([]byte, packetSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(packet)
		if err != nil || n != 3 {
			slog.Warn(fmt.Sprintf("invalid packet from %s", remoteAddr), "handleListener", err.Error())
			continue
		}

		protocol, port := parsePacket(packet)
		slog.Debug(fmt.Sprintf("received %s %d from %s", protocol, port, remoteAddr.IP.To16()))

		if !checkPort(port) || !checkIP(remoteAddr.IP.To16().String()) {
			slog.Debug("denied ip/port", "port", port, "remoteAddr", remoteAddr)
			continue
		}

		allocations := GetchAllocations()
		if allocations == nil {
			continue
		}

		if checkPortAllocation(conn, allocations, port, remoteAddr) {
			continue
		}

		forwardPort(conn, remoteAddr, port, protocol)
	}
}

func parsePacket(packet []byte) (types.Protocol, uint16) {
	protocol := types.TCP
	if string(packet[:1]) == "u" {
		protocol = types.UDP
	}

	port := binary.BigEndian.Uint16(packet[1:])
	return protocol, port
}

func checkIP(ip string) bool {
	slog.Debug("checking ip", "ip", ip)
	return len(config.Server.AllowedIPs) == 0 || slices.Contains(config.Server.AllowedIPs, ip)
}

func checkPort(port uint16) bool {
	slog.Debug("checking port", "port", port)
	return len(config.Server.AllowedPorts) == 0 || slices.Contains(config.Server.AllowedPorts, port)
}

func checkPortAllocation(conn *net.UDPConn, allocations []types.Allocation, port uint16, remoteAddr *net.UDPAddr) bool {
	for _, alloc := range allocations {
		slog.Debug(fmt.Sprintf("checking port: %d\n", alloc.Port))
		if alloc.Port == port {
			conn.WriteToUDP([]byte{byte(types.Allocated)}, remoteAddr)
			slog.Debug(fmt.Sprintf("port %d is already allocated\n", port))
			return true
		}
	}
	return false
}

func forwardPort(conn *net.UDPConn, remoteAddr *net.UDPAddr, port uint16, protocol types.Protocol) {
	err := utils.ForwardPort(remoteAddr.IP.To16().String(), uint16(port), protocol)
	if err != nil {
		slog.Error("failed to forward port", "forwardPort", err.Error())
		return
	}

	slog.Info(fmt.Sprintf("port forwarded: %d/%s -> %s\n", port, protocol, remoteAddr.IP.To16()))

	_, err = conn.WriteToUDP([]byte{byte(types.OK)}, remoteAddr)
	if err != nil {
		slog.Error("failed to reply", "forwardPort", err.Error())
	}
}
