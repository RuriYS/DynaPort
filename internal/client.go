package internal

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/RuriYS/DynaPort/types"
)

func StartClient(config types.Config) {
	slog.Info("dynaport is alive!")

	serverAddr := net.UDPAddr{
		IP:   net.ParseIP(config.Client.Host),
		Port: int(config.Client.Port),
	}

	broadcastPeriod, err := time.ParseDuration(config.Client.BroadcastInterval)
	if err != nil {
		slog.Error("failed to parse broadcast_interval", "StartClient", err.Error())
	}

	timeout, err := time.ParseDuration(config.Client.Timeout)
	if err != nil {
		slog.Error("failed to parse timeout", "StartClient", err.Error())
	}

	for {
		allocations := GetchAllocations()
		conn := establishConnection(&serverAddr, broadcastPeriod)
		if conn == nil {
			continue
		}

		processAllocations(conn, allocations, timeout)

		conn.Close()
		slog.Info(fmt.Sprintf("broadcast complete, sleeping for %v", broadcastPeriod))
		time.Sleep(broadcastPeriod)
	}
}

func establishConnection(serverAddr *net.UDPAddr, broadcastPeriod time.Duration) *net.UDPConn {
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		slog.Error("failed to connect", "establishConnection", err.Error())
		time.Sleep(broadcastPeriod)
		return nil
	}
	slog.Debug(fmt.Sprintf("connected to %v", serverAddr))
	return conn
}

func processAllocations(conn *net.UDPConn, allocations []types.Allocation, timeout time.Duration) {
	for _, alloc := range allocations {
		var protoByte byte
		if alloc.Protocol == types.TCP {
			protoByte = 't'
		} else {
			protoByte = 'u'
		}

		packet := make([]byte, packetSize)
		packet[0] = protoByte
		binary.BigEndian.PutUint16(packet[1:], alloc.Port)
		_, err := conn.Write(packet)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to send packet for %s %d", alloc.Protocol, alloc.Port), "processAllocations", err.Error())
			continue
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		resp := make([]byte, 1)
		n, _, err := conn.ReadFromUDP(resp)
		if err != nil {
			slog.Warn(fmt.Sprintf("no response for %s %d", alloc.Protocol, alloc.Port), "processAllocations", err.Error())
			continue
		}

		res := resp[:n]
		if len(res) == 1 && res[0] == byte(types.OK) {
			slog.Info(fmt.Sprintf("forwarded port %s %d", alloc.Protocol, alloc.Port))
		} else if len(res) == 1 && res[0] == byte(types.Allocated) {
			slog.Warn(fmt.Sprintf("port already allocated %s %d", alloc.Protocol, alloc.Port))
		}
	}
}
