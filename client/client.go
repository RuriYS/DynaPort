package client

import "log/slog"

func Run() {
	slog.Info("[Client] dynaport is alive!")
	StartRelay()
}
