package client

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/RuriYS/RePorter/internal"
	"github.com/RuriYS/RePorter/types"
)

func StartRelay() {
	config := internal.GetConfig()

	serverAddr := net.UDPAddr{
		IP:   net.ParseIP(config.Client.Host),
		Port: int(config.Client.Port),
	}

	interval, err := time.ParseDuration(config.Client.BroadcastInterval)
	if err != nil {
		slog.Error("[Relay] failed to parse broadcast_interval", "error", err.Error(), "broadcast_interval", config.Client.BroadcastInterval)
	}

	timeout, err := time.ParseDuration(config.Client.Timeout)
	if err != nil {
		slog.Error("[Relay] failed to parse timeout", "error", err.Error(), "timeout", config.Client.Timeout)
	}

	for {
		sockets := internal.GetAllocations()
		conn := connect(&serverAddr, interval)
		if conn == nil {
			continue
		}

		for _, socket := range sockets {
			send(conn, &socket, timeout)
		}

		conn.Close()
		slog.Info(fmt.Sprintf("[Relay] broadcast complete, sleeping for %v", interval))
		time.Sleep(interval)
	}
}

func connect(serverAddr *net.UDPAddr, broadcastPeriod time.Duration) *net.UDPConn {
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		slog.Error("[Relay::connect] failed to connect", "error", err.Error())
		time.Sleep(broadcastPeriod)
		return nil
	}
	slog.Debug("[Relay::connect] connected", "serverAddr", serverAddr, "broadcastPeriod", broadcastPeriod)
	return conn
}

func send(conn *net.UDPConn, socket *types.Allocation, timeout time.Duration) {
	var protoByte byte
		if socket.Protocol == types.TCP {
			protoByte = 't'
		} else {
			protoByte = 'u'
		}

		packet := make([]byte, 3)
		packet[0] = protoByte
		binary.BigEndian.PutUint16(packet[1:], socket.Port)
		_, err := conn.Write(packet)
		if err != nil {
			slog.Error("[Relay::send] failed to send packet", "error", err.Error(), "packet", packet)
			return
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		resp := make([]byte, 1)
		n, _, err := conn.ReadFromUDP(resp)
		if err != nil {
			slog.Warn("[Relay::send] timed out", "error", err.Error())
			return
		}

		res := resp[:n]
		if len(res) == 1 && res[0] == byte(types.OK) {
			slog.Info("[Relay::send] port forwarded", "port", socket.Port, "protocol", socket.Protocol)
		} else if len(res) == 1 && res[0] == byte(types.Allocated) {
			slog.Debug("[Relay::send] port already allocated", "port", socket.Port, "protocol", socket.Protocol)
		}
}
