package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/RuriYS/DynaPorts/utils"
	"golang.org/x/exp/slog"
)

func main()  {
	slog.Info("dynaport is alive!")

	addr := net.UDPAddr{Port: 10000, IP: net.ParseIP("0.0.0.0")}

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        slog.Error("%s", err.Error())
    }

    defer conn.Close()

    buffer := make([]byte, 6)
    for {
        n, remoteAddr, err := conn.ReadFromUDP(buffer)
        if n < 2 {
            slog.Warn(fmt.Sprintf("Received malformed packet from %s", remoteAddr))
            continue
        }
        if err != nil {
            slog.Error("%s", err.Error())
        }

        prot_str := string(buffer[:1])
        protocol := utils.TCP
        if prot_str == "u" {
            protocol = utils.UDP
        }
        port, err := strconv.ParseUint(strings.Trim(string(buffer[1:n]), "\n"), 10, 64)
        if err != nil {
            slog.Error(err.Error())
            continue
        }
        
        slog.Info(fmt.Sprintf("Received %s %d from %s", protocol, port, remoteAddr.IP.To16()))
        err = utils.ForwardPort(remoteAddr.IP.To16().String(), uint16(port), protocol)
        if err != nil {
            slog.Error(err.Error())
        }
    }
}
