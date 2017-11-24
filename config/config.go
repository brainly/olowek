package config

import (
	"encoding/json"
	"io/ioutil"
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
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
