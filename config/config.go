package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/brainly/olowek/marathon"
)

const (
	DefaultNginxConfig   = "/etc/nginx/conf.d/services.conf"
	DefaultNginxTemplate = "/etc/olowek/services.tpl"
	DefaultNginxCmd      = "/usr/sbin/nginx"
	EmptyScope           = ""
)

type Config struct {
	sync.RWMutex
	Scope           string `json:"scope,omitempty"`
	Marathon        string `json:"marathon"`
	NginxConfig     string `json:"nginx_config,omitempty"`
	NginxTemplate   string `json:"nginx_template,omitempty"`
	NginxCmd        string `json:"nginx_cmd,omitempty"`
	NginxReloadFunc func(string) error
	Apps            []marathon.Application
}

func NewConfigFromFile(path string) (*Config, error) {
	config := Config{
		Scope:         EmptyScope,
		NginxConfig:   DefaultNginxConfig,
		NginxTemplate: DefaultNginxTemplate,
		NginxCmd:      DefaultNginxCmd,
	}

	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return nil, err
	}

	if config.Marathon == "" {
		return nil, fmt.Errorf("Missing 'marathon' field in '%s'", path)
	}

	return &config, nil
}
