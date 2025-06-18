package client

import (
	"log/slog"
)


func Run() {
	slog.Info("[Client] starting")
	StartRelay()
}
