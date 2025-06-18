package server

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"

	"github.com/RuriYS/RePorter/internal/firewalld"
	"github.com/RuriYS/RePorter/internal/sockit"
	"github.com/RuriYS/RePorter/types"
)

func StartListener(conn *net.UDPConn) {
	slog.Info("[Listener] listening", "addr", conn.LocalAddr().String())
	buffer := make([]byte, 3)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil || n != 3 {
			slog.Warn("[Listener] invalid packet", "remoteAddr", remoteAddr, "error", err.Error())
			continue
		}

		protocol, port := parsePacket(buffer)
		slog.Debug("[Listener] received message", "remoteAddr", remoteAddr, "protocol", protocol, "port", port)

		if !checkPort(port) || !checkIP(remoteAddr.IP.String()) {
			slog.Debug("[Listener] denied ip/port", "port", port, "remoteAddr", remoteAddr)
			continue
		}

		allocations := sockit.GetAll()
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
	if packet[0] == 'u' {
		protocol = types.UDP
	}
	port := binary.BigEndian.Uint16(packet[1:])
	return protocol, port
}

func checkIP(ip string) bool {
    slog.Debug("[checkIP] checking ip", "ip", ip)
    return len(allowedIPs) == 0 || hasIP(ip)
}

func hasIP(ip string) bool {
    _, ok := allowedIPs[ip]
    return ok
}

func checkPort(port uint16) bool {
    slog.Debug("[checkPort] checking port", "port", port)
    return len(allowedPorts) == 0 || hasPort(port)
}

func hasPort(port uint16) bool {
    _, ok := allowedPorts[port]
    return ok
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
	err := firewalld.ForwardPort(remoteAddr.IP.String(), uint16(port), protocol)
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
