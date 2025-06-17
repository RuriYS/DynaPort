package internal

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/RuriYS/DynaPort/types"
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
  host: 127.0.0.1
  port: 42000
  broadcast_interval: 5m
  timeout: 3s
  ports:
   - 22
  whitelist_mode: false # forward ports in the list if true, otherwise don't (default: false)
`

func GetConfig(path string) (config types.Config, err error) {
	if len(path) == 0 {
		path = "/etc/dynaport/config.yml"
	}

	slog.Debug("getting config", "GetConfig:path", path)
	
	_, err = os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Debug("config does not exist")
			createConfig(path)
		}
	}
	
	slog.Debug("reading config", "GetConfig:path", path)
	f, err := os.ReadFile(path)
	if err != nil {
		return types.Config{}, nil
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return types.Config{}, nil
	}

	slog.Debug("found config", "GetConfig", config)
	return config, nil
}

func createConfig(path string) {
	slog.Debug("creating directory", "createConfig:path", path)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		slog.Error("failed to create directory", "GetConfig", err.Error())
		os.Exit(1)
	}
	slog.Debug("writing default config", "createConfig:path", path)
	err = os.WriteFile(path, []byte(default_config), 0755)
	if err != nil {
		slog.Error("failed to create config", "GetConfig", err.Error())
		os.Exit(1)
	}
}
