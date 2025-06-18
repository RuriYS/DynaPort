package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/RuriYS/RePorter/internal/config"
	"github.com/RuriYS/RePorter/internal/sockit"
	"github.com/RuriYS/RePorter/types"
)

func StartRelay() {
	config := config.GetConfig()
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
		sockets := sockit.GetAll()
		conn := connect(&serverAddr, interval)
		if conn == nil {
			slog.Error("[Relay] failed to connect", "error", err.Error())
			continue
		}
		slog.Debug("[Relay] connected", "serverAddr", serverAddr, "broadcastPeriod", interval)

		for _, socket := range sockets {
			slog.Debug("[Relay] broadcasting socket", "socket", socket)
			packet, err := send(conn, &socket, timeout)
			if err != nil {
				slog.Error("[Relay] broadcast failed", "error", err.Error(), "packet", packet)
				continue
			}
			if packet[0] == byte(types.OK) {
				slog.Info("[Relay] socket broadcasted", "socket", socket)
			} else if packet[0] == byte(types.Allocated) {
				slog.Warn("[Relay] port already allocated", "socket", socket)
			}
			
			slog.Info("[Relay] socket broadcasted", "socket", socket)
		}

		conn.Close()
		slog.Info(fmt.Sprintf("[Relay] broadcast complete, sleeping for %v", interval))
		time.Sleep(interval)
	}
}

func connect(serverAddr *net.UDPAddr, interval time.Duration) *net.UDPConn {
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		time.Sleep(interval)
		return nil
	}
	return conn
}

func send(conn *net.UDPConn, socket *types.Allocation, timeout time.Duration) (packet []byte, err error) {
	var protoByte byte
	if socket.Protocol == types.TCP {
		protoByte = 't'
	} else {
		protoByte = 'u'
	}

	packet = make([]byte, 3)
	packet[0] = protoByte
	binary.BigEndian.PutUint16(packet[1:], socket.Port)
	_, err = conn.Write(packet)
	if err != nil {
		return packet, err
	}

	conn.SetReadDeadline(time.Now().Add(timeout))
	resp := make([]byte, 1)
	n, _, err := conn.ReadFromUDP(resp)
	if err != nil {
		return packet, err
	}

	res := resp[:n]
	if len(res) == 1 {
		return res, nil
	}

	return res, errors.New("received invalid packet")
}
