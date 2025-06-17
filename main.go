package main

import (
	"flag"
	"log"
	"log/slog"

	"github.com/RuriYS/DynaPort/internal"
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
        log.Fatalln("ERROR: --server and --client cannot both be set")
    }

    if *verbose {
        slog.SetLogLoggerLevel(slog.LevelDebug)
    }

    config, err := internal.GetConfig(*configPath)
    if err != nil {
        slog.Error("failed to load config", "main", err.Error())
    }

	if *serverMode {
		internal.StartServer(config)
	} else if *clientMode {
		internal.StartClient(config)
	}
}
