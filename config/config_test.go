package config

import "testing"

func TestNewConfigFromFile(t *testing.T) {
	configCases := map[string]struct {
		fixture     string
		config      *Config
		expectedErr bool
	}{
		"Test valid configuration": {
			fixture: "./fixtures/olowek.json",
			config: &Config{
				Scope:         "internal",
				Marathon:      "http://127.0.0.1:8080,127.0.0.1:8080",
				NginxConfig:   "services.conf",
				NginxTemplate: "services.tpl",
				NginxCmd:      "/usr/sbin/nginx",
			},
		},
	}
	for name, tt := range configCases {
		t.Run(name, func(t *testing.T) {
			c, err := NewConfigFromFile(tt.fixture)

			if tt.expectedErr && err != nil {
				t.Fatalf("Unexpected error: '%s'", err)
			}

			if c.Scope != tt.config.Scope ||
				c.Marathon != tt.config.Marathon ||
				c.NginxConfig != tt.config.NginxConfig ||
				c.NginxTemplate != tt.config.NginxTemplate ||
				c.NginxCmd != tt.config.NginxCmd {
				t.Fatalf("Configs are not equal. Got: %#v", c)
			}
		})
	}

}
