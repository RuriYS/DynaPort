package utils

import (
	"log/slog"

	"github.com/RuriYS/DynaPorts/types"
	"github.com/cakturk/go-netstat/netstat"
)

type Allocation struct {
    Protocol types.Protocol
    Port uint16
}

func GetAllocations() []Allocation {
    var allocations []Allocation

	socks, err := netstat.TCPSocks(func(ste *netstat.SockTabEntry) bool {
        return ste.State == netstat.Listen
    })
	if err != nil {
		slog.Error(err.Error())
	}

	for _, sock := range socks {
        allocations = append(allocations, Allocation{
            Protocol: types.TCP,
            Port: sock.RemoteAddr.Port,
        })
	}

	socks, err = netstat.UDPSocks(netstat.NoopFilter)
	if err != nil {
		slog.Error(err.Error())
	}

	for _, sock := range socks {
		allocations = append(allocations, Allocation{
            Protocol: types.UDP,
            Port: sock.RemoteAddr.Port,
        })
	}

	return allocations
}
