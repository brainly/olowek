package config

import (
	"encoding/json"
	"os"
	"sync"
)

type App struct {
	Name   string
	ID     string
	Labels map[string]string
	Env    map[string]string
	Tasks  []AppTask
}

type AppTask struct {
	ID           string
	Host         string
	Ports        []int
	ServicePorts []int
}

type Config struct {
	sync.RWMutex
	Scope           string `json:"scope"`
	Marathon        string `json:"marathon"`
	NginxConfig     string `json:"nginx_config"`
	NginxTemplate   string `json:"nginx_template"`
	NginxCmd        string `json:"nginx_cmd"`
	NginxReloadFunc func(string) error
	Apps            []App
}

func NewConfigFromFile(path string) (*Config, error) {
	var config Config

	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	return &config, err
}
