package config

import (
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
	Scope           string
	Marathon        string
	NginxConfig     string
	NginxTemplate   string
	NginxCmd        string
	NginxReloadFunc func(string) error
	Apps            []App
}
