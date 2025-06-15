package utils

import (
	"strconv"
	"strings"
)

func ParsePort(data []byte) (uint16, error) {
	portStr := strings.TrimSpace(string(data))
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(port), nil
}
