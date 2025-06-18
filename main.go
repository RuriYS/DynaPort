package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/RuriYS/RePorter/client"
	"github.com/RuriYS/RePorter/internal"
	"github.com/RuriYS/RePorter/server"
)

func main() {
    serverMode := flag.Bool("server", false, "Run DynaPort as server mode")
    flag.BoolVar(serverMode, "s", false, "Alias for --server")

    clientMode := flag.Bool("client", false, "Run DynaPort as client mode")
    flag.BoolVar(clientMode, "c", false, "Alias for --client")

    configPath := flag.String("config", "", "Config path")
    
    verbose := flag.Bool("verbose", false, "Verbose logging")
    flag.BoolVar(verbose, "v", false, "Alias for --verbose")

    flag.Parse()

    if *serverMode && *clientMode {
        slog.Error("[main] server and client cannot both be set")
        os.Exit(1)
    }

    if *verbose {
        slog.SetLogLoggerLevel(slog.LevelDebug)
    }

    err := internal.LoadConfig(*configPath)
    if err != nil {
        slog.Error("[main] failed to load config", "error", err.Error())
        os.Exit(1)
    }

    mode := "server"
    if *clientMode {
        mode = "client"
    } else {
        *serverMode = true
    }
    
    slog.Info("[main] Starting " + mode)

    internal.RunAllocator()
    
	if *serverMode {
		server.Run()
	} else if *clientMode {
		client.Run()
	}
}
