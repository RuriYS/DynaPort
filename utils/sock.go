package utils

import (
	"github.com/RuriYS/DynaPort/types"
	"github.com/cakturk/go-netstat/netstat"
)

func GetAllocations() (a []types.Allocation, err error) {
	var allocations []types.Allocation

	socks, err := netstat.TCPSocks(func(ste *netstat.SockTabEntry) bool {
		return ste.State == netstat.Listen
	})
	if err != nil {
		return nil, err
	}

	for _, sock := range socks {
		allocations = append(allocations, types.Allocation{
			Protocol: types.TCP,
			Port:     sock.LocalAddr.Port,
		})
	}

	socks, err = netstat.UDPSocks(netstat.NoopFilter)
	if err != nil {
		return nil, err
	}

	for _, sock := range socks {
		allocations = append(allocations, types.Allocation{
			Protocol: types.UDP,
			Port:     sock.LocalAddr.Port,
		})
	}

	return allocations, nil
}
