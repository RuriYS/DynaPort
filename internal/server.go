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

	slog.Info("[StartServer] starting listener", "addr", addr)
	handleListener(conn)
}

func initializeServer(addr *net.UDPAddr) *net.UDPConn {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("[initializeServer] failed to start the server", "error", err.Error())
		return nil
	}

	slog.Info("[initializeServer] dynaport is alive!")
	return conn
}

func handleListener(conn *net.UDPConn) {
	packet := make([]byte, packetSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(packet)
		if err != nil || n != 3 {
			slog.Warn("[handleListener] invalid packet", "remoteAddr", remoteAddr, "error", err.Error())
			continue
		}

		protocol, port := parsePacket(packet)
		slog.Debug("[handleListener] received message", "remoteAddr", remoteAddr, "protocol", protocol, "port", port)

		if !checkPort(port) || !checkIP(remoteAddr.IP.To16().String()) {
			slog.Debug("[handleListener] denied ip/port", "port", port, "remoteAddr", remoteAddr)
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
	slog.Debug("[checkIP] checking ip", "ip", ip)
	return len(config.Server.AllowedIPs) == 0 || slices.Contains(config.Server.AllowedIPs, ip)
}

func checkPort(port uint16) bool {
	slog.Debug("[checkPort] checking port", "port", port)
	return len(config.Server.AllowedPorts) == 0 || slices.Contains(config.Server.AllowedPorts, port)
}

func checkPortAllocation(conn *net.UDPConn, allocations []types.Allocation, port uint16, remoteAddr *net.UDPAddr) bool {
	for _, alloc := range allocations {
		slog.Debug("[checkPortAllocation] checking port", "port", alloc.Port)
		if alloc.Port == port {
			conn.WriteToUDP([]byte{byte(types.Allocated)}, remoteAddr)
			slog.Debug("[checkPortAllocation] port already allocated", "remoteAddr", remoteAddr, "port", port, "allocations", allocations)
			return true
		}
	}
	return false
}

func forwardPort(conn *net.UDPConn, remoteAddr *net.UDPAddr, port uint16, protocol types.Protocol) {
	err := utils.ForwardPort(remoteAddr.IP.To16().String(), uint16(port), protocol)
	if err != nil {
		slog.Error("[forwardPort] failed to forward port", "remoteAddr", remoteAddr, "port", port, "protocol", protocol, "error", err.Error())
		return
	}

	slog.Info(fmt.Sprintf("[forwardPort] port forwarded: %d/%s -> %s\n", port, protocol, remoteAddr.IP))

	_, err = conn.WriteToUDP([]byte{byte(types.OK)}, remoteAddr)
	if err != nil {
		slog.Error("[forwardPort] failed to reply", "error", err.Error())
	}
}
