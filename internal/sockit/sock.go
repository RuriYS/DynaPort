package sockit

import (
	"log/slog"
	"sync"
	"time"

	"github.com/RuriYS/RePorter/internal/config"
	"github.com/RuriYS/RePorter/types"
)

var (
    cache		[]types.Allocation
    mutex		sync.RWMutex
	interval	time.Duration
	ready		chan struct{}
)

func initialize() {
	config := config.GetConfig()
	var err error
	interval, err = time.ParseDuration(config.Client.BroadcastInterval)
	if err != nil {
		slog.Error("[Sockit] initialization failed", "error", err)
		return
	}
	cache, err = getSockets()
	if err != nil {
		slog.Error("[Sockit] failed to get initial sockets", "error", err)
		return
	}
	if len(cache) == 0 {
		slog.Warn("[Sockit] no sockets found, will retry in the next interval")
		return
	}
}

func Run() chan struct{} {
	ready = make(chan struct{})
	go func ()  {
		slog.Info("[Sockit] initializing")
		initialize()
		close(ready)
		
		tickDuration := interval - time.Second
		if tickDuration < time.Second {
			tickDuration = time.Second
			slog.Warn("[Sockit] interval too short, using minimum 1s ticker")
		}
		ticker := time.NewTicker(tickDuration)
		
		defer ticker.Stop()
		for {
			<- ticker.C
			slog.Debug("[Sockit] scanning for sockets")
			allocs, err := getSockets()
			if err != nil {
				slog.Error("failed to get allocations", "error", err.Error())
				continue
			}
			slog.Debug("[Sockit] sockets found", "allocs", allocs)
			mutex.Lock()
			cache = allocs
			mutex.Unlock()
		}
	}()
	return ready
}

func GetAll() []types.Allocation {
	mutex.RLock()
	defer mutex.RUnlock()
	return cache
}

func getSockets() (allocs []types.Allocation, err error) {
	allocs, err = GetSocks()
	if err != nil {
		return nil, err
	}
	return allocs, nil
}
