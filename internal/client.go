package internal

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/RuriYS/DynaPort/types"
)

func StartClient() {
	slog.Info("[StartClient] dynaport is alive!")
	handleBroadcast()
}

func handleBroadcast() {
	serverAddr := net.UDPAddr{
		IP:   net.ParseIP(config.Client.Host),
		Port: int(config.Client.Port),
	}

	broadcastPeriod, err := time.ParseDuration(config.Client.BroadcastInterval)
	if err != nil {
		slog.Error("[handleBroadcast] failed to parse broadcast_interval", "error", err.Error(), "broadcast_interval", config.Client.BroadcastInterval)
	}

	timeout, err := time.ParseDuration(config.Client.Timeout)
	if err != nil {
		slog.Error("[handleBroadcast] failed to parse timeout", "error", err.Error(), "timeout", config.Client.Timeout)
	}

	for {
		allocations := GetAllocations()
		conn := establishConnection(&serverAddr, broadcastPeriod)
		if conn == nil {
			continue
		}

		processAllocations(conn, allocations, timeout)

		conn.Close()
		slog.Info(fmt.Sprintf("[handleBroadcast] broadcast complete, sleeping for %v", broadcastPeriod))
		time.Sleep(broadcastPeriod)
	}
}

func establishConnection(serverAddr *net.UDPAddr, broadcastPeriod time.Duration) *net.UDPConn {
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		slog.Error("[establishConnection] failed to connect", "error", err.Error())
		time.Sleep(broadcastPeriod)
		return nil
	}
	slog.Debug("[establishConnection] connected", "serverAddr", serverAddr, "broadcastPeriod", broadcastPeriod)
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
			slog.Error("[processAllocations] failed to send packet", "error", err.Error(), "packet", packet)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		resp := make([]byte, 1)
		n, _, err := conn.ReadFromUDP(resp)
		if err != nil {
			slog.Warn("[processAllocations] timed out", "error", err.Error())
			continue
		}

		res := resp[:n]
		if len(res) == 1 && res[0] == byte(types.OK) {
			slog.Info("[processAllocations] port forwarded", "port", alloc.Port, "protocol", alloc.Protocol)
		} else if len(res) == 1 && res[0] == byte(types.Allocated) {
			slog.Debug("[processAllocations] port already allocated", "port", alloc.Port, "protocol", alloc.Protocol)
		}
	}
}
