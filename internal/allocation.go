package internal

import (
	"log/slog"

	"github.com/RuriYS/DynaPort/types"
	"github.com/RuriYS/DynaPort/utils"
)

func GetAllocations() []types.Allocation {
	allocs, err := utils.GetAllocations()
	if err != nil {
		slog.Error("failed to get allocations", "error", err.Error())
		return nil
	}
	return allocs
}

