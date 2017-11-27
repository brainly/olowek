package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/brainly/olowek/marathon"
)

var (
	ErrMissingMarathon      = errors.New("Missing 'marathon' filed in configuration")
	ErrMissingNginxConfig   = errors.New("Missing 'nginx_config' filed in configuration")
	ErrMissingNginxTemplate = errors.New("Missing 'nginx_template' filed in configuration")
	ErrMissingNginxCmd      = errors.New("Missing 'nginx_cmd' filed in configuration")
)

type Config struct {
	sync.RWMutex
	Scope           string `json:"scope"`
	Marathon        string `json:"marathon"`
	NginxConfig     string `json:"nginx_config"`
	NginxTemplate   string `json:"nginx_template"`
	NginxCmd        string `json:"nginx_cmd"`
	NginxReloadFunc func(string) error
	Apps            []marathon.Application
}

func NewConfigFromFile(path string) (*Config, error) {
	var config Config

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
		return nil, ErrMissingMarathon
	}

	if config.NginxConfig == "" {
		return nil, ErrMissingNginxConfig
	}

	if config.NginxTemplate == "" {
		return nil, ErrMissingNginxTemplate
	}

	if config.NginxCmd == "" {
		return nil, ErrMissingNginxCmd
	}

	return &config, nil
}
