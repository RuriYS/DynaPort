package internal

import (
	"log/slog"
	"sync"
	"time"

	"github.com/RuriYS/RePorter/types"
	"github.com/RuriYS/RePorter/utils"
)

var (
    cache		[]types.Allocation
    mutex		sync.RWMutex
	interval	time.Duration
	ready		chan struct{}
)

func initialize() {
	config := GetConfig()
	var err error
	interval, err = time.ParseDuration(config.Client.BroadcastInterval)
	if err != nil {
		slog.Error("[Allocator] initialization failed", "error", err)
		return
	}
}

func RunAllocator() chan struct{} {
	ready = make(chan struct{})
	go func ()  {
		slog.Info("[Allocator] initializing")
		initialize()
		cache = getSockets()
		slog.Debug("[Allocator] sockets cached", "cache", cache)
		close(ready)
		ticker := time.NewTicker(interval)
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
	return ready
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
