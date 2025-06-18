package internal

import (
	"log/slog"
	"sync"
	"time"

	"github.com/RuriYS/RePorter/types"
	"github.com/RuriYS/RePorter/utils"
)

var (
    cache	[]types.Allocation
    mutex	sync.RWMutex
)

func RunAllocator() {
	go func ()  {
		slog.Debug("[Allocator] initialized")
		cache = getSockets()
		slog.Debug("[Allocator] sockets cached", "cache", cache)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			<- ticker.C
			slog.Debug("[Allocator] scanning for sockets")
			allocs := getSockets()
			slog.Debug("[Allocator] sockets found", "allocs", allocs)
			mutex.Lock()
			cache = allocs
			mutex.Unlock()
		}
	}()
}

func GetAllocations() []types.Allocation {
	mutex.RLock()
	defer mutex.RUnlock()
	return cache
}

func getSockets() []types.Allocation {
	allocs, err := utils.GetAllocations()
	if err != nil {
		slog.Error("failed to get allocations", "error", err.Error())
		return nil
	}
	return allocs
}
