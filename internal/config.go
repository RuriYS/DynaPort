package internal

import (
	"errors"
	"go/types"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const default_config = `
server:
  host: 0.0.0.0
  port: 42000
  ttl: 3h
  allowed_ips: 
   - 0.0.0.0
  allowed_ports:
   - 8080
   - 8443

client:
  host: 127.0.0.1:42000
  broadcast_interval: 5m
  timeout: 3s
  ports:
   - 22
  whitelist_mode: false # forward ports in the list if true, otherwise don't (default: false)
`

func GetConfig(path string) (config types.Config, err error) {
	if len(path) == 0 {
		path = "/etc/dynaport/config"
	}

	_, err = os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			createConfig(path)
		}
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return types.Config{}, nil
	}

	c := &types.Config{}
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return types.Config{}, nil
	}

	return *c, nil
}

func createConfig(path string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		slog.Error("failed to create directory", "GetConfig", err.Error())
		os.Exit(1)
	}
	err = os.WriteFile(path, []byte(strings.TrimSpace(default_config)), 0600)
	if err != nil {
		slog.Error("failed to create config", "GetConfig", err.Error())
		os.Exit(1)
	}
}
