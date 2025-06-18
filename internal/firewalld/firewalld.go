package firewalld

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/RuriYS/RePorter/types"
)

func ForwardPort(addr string, port uint16, proto types.Protocol) error {
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

func UnforwardPort(addr string, port uint16, proto types.Protocol) error {
	cmd := exec.Command("firewall-cmd", fmt.Sprintf("--remove-forward-port=toaddr=%s:port=%d:proto=%s", addr, port, proto))
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
