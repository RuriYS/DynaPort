package internal

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/RuriYS/DynaPort/types"
	"github.com/RuriYS/DynaPort/utils"
)

const (
	timeoutSec      = 5
	broadcastPeriod = 3 * time.Minute
)

func StartClient(serverHost string, serverPort uint16) {
    serverAddr := net.UDPAddr{
        IP:   net.ParseIP(serverHost),
        Port: int(serverPort),
    }

    for {
        allocs, err := utils.GetAllocations()
        if err != nil {
            slog.Error(fmt.Sprintf("failed to get allocations: %v", err))
            time.Sleep(broadcastPeriod)
            continue
        }

        conn, err := net.DialUDP("udp", nil, &serverAddr)
        if err != nil {
            slog.Error(fmt.Sprintf("failed to dial server: %v", err))
            time.Sleep(broadcastPeriod)
            continue
        }

        for _, alloc := range allocs {
            var protoByte byte
            if alloc.Protocol == types.TCP {
                protoByte = 't'
            } else {
                protoByte = 'u'
            }

            packet := append([]byte{protoByte}, byte(alloc.Port))

            _, err := conn.Write(packet)
            if err != nil {
                slog.Error(fmt.Sprintf("failed to send packet for %s %d: %v", alloc.Protocol, alloc.Port, err))
                continue
            }

            conn.SetReadDeadline(time.Now().Add(timeoutSec * time.Second))
            resp := make([]byte, 1)
            n, _, err := conn.ReadFromUDP(resp)
            if err != nil {
                slog.Warn(fmt.Sprintf("no response for %s %d: %v", alloc.Protocol, alloc.Port, err))
                continue
            }

            res := resp[:n]
            if len(res) == 1 && res[0] == byte(OK) {
                slog.Info(fmt.Sprintf("forwarded port %s %d", alloc.Protocol, alloc.Port))
            } else if len(res) == 1 && res[0] == byte(Allocated) {
                slog.Warn(fmt.Sprintf("port already allocated %s %d", alloc.Protocol, alloc.Port))
            }
        }

        conn.Close()
        slog.Info(fmt.Sprintf("broadcast complete, sleeping for %v", broadcastPeriod))
        time.Sleep(broadcastPeriod)
    }
}
