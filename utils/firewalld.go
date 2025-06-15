package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Protocol string

const (
    TCP Protocol = "tcp"
    UDP Protocol = "udp"
)

func ForwardPort(addr string, port uint16, proto Protocol) error {
    allocations := GetAllocations()
    for _, alloc := range allocations {
        if alloc.Port == port {
            return fmt.Errorf("port %d is already allocated", port)
        }
    }

	cmd := exec.Command("firewall-cmd", fmt.Sprintf("--add-forward-port=toaddr=%s:port=%d:proto=%s", addr, port, proto))
    o, err := cmd.Output()
    output := string(o)
    if err != nil {
        return err
    }
    if !strings.Contains(output, "success") {
        return errors.New(output)
    }
    
    return nil
}
