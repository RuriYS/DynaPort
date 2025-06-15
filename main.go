package main

import (
	"flag"
	"log"

	"github.com/RuriYS/DynaPort/internal"
)

func main() {
    host := flag.String("host", "0.0.0.0", "Host for DynaPort server (default: 0.0.0.0)")
    flag.StringVar(host, "h", "0.0.0.0", "Alias for --host")

    port := flag.Uint("port", 10000, "Port for DynaPort server (default: 10000)")
    flag.UintVar(port, "p", 10000, "Alias for --port")

    serverMode := flag.Bool("server", false, "Run DynaPort as server mode")
    flag.BoolVar(serverMode, "s", false, "Alias for --server")

    clientMode := flag.Bool("client", false, "Run DynaPort as client mode")
    flag.BoolVar(clientMode, "c", false, "Alias for --client")
    
    verbose := flag.Bool("verbose", false, "Verbose logging")
    flag.BoolVar(verbose, "v", false, "Alias for --verbose")

    flag.Parse()

    if *serverMode && *clientMode {
        log.Fatalln("ERROR: --server and --client cannot both be set")
    }

	if *serverMode {
		internal.StartServer(*host, uint16(*port), *verbose)
	} else if *clientMode {
		internal.StartClient(*host, uint16(*port), *verbose)
	}
}
