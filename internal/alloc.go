package internal

import (
	"log/slog"
	"sync"
	"time"

	"github.com/RuriYS/DynaPort/types"
	"github.com/RuriYS/DynaPort/utils"
)

var (
    cache	[]types.Allocation
    mutex	sync.RWMutex
)

func RunAllocator() {
	go func ()  {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			<- ticker.C
			allocs := getSockets()
			mutex.Lock()
			cache = allocs
			mutex.Lock()
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
