package config

import (
	"github.com/go-test/deep"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	configCases := map[string]struct {
		fixture    string
		config     *Config
		shouldFail bool
	}{
		"Configuration with all fileds set properly": {
			fixture:    "./fixtures/olowek.json",
			shouldFail: false,
			config: &Config{
				Scope:         "internal",
				Marathon:      "http://127.0.0.1:8080,127.0.0.1:8080",
				NginxConfig:   "services.conf",
				NginxTemplate: "services.tpl",
				NginxCmd:      "/usr/sbin/nginx",
			},
		},
		"Missing scope should produce empty scope": {
			fixture:    "./fixtures/noscope.json",
			shouldFail: false,
			config: &Config{
				Scope:         "",
				Marathon:      "http://127.0.0.1:8080,127.0.0.1:8080",
				NginxConfig:   "services.conf",
				NginxTemplate: "services.tpl",
				NginxCmd:      "/usr/sbin/nginx",
			},
		},
		"Nonexisting config should fail and return config set to nil": {
			fixture:    "i_do_not_exist.json",
			shouldFail: true,
			config:     nil,
		},
		"Missing marathon should fail": {
			fixture:    "./fixtures/missing_marathon.json",
			shouldFail: true,
			config:     nil,
		},
		"Can specify only marathon field, other should be set to default values": {
			fixture:    "./fixtures/only_marathon.json",
			shouldFail: false,
			config: &Config{
				Scope:         EmptyScope,
				NginxConfig:   DefaultNginxConfig,
				NginxTemplate: DefaultNginxTemplate,
				NginxCmd:      DefaultNginxCmd,
				Marathon:      "http://marathon:8080",
			},
		},
	}
	for name, tt := range configCases {
		t.Run(name, func(t *testing.T) {
			c, err := NewConfigFromFile(tt.fixture)

			if !tt.shouldFail && err != nil {
				t.Fatalf("Unexpected error: '%s'", err)
			}

			if tt.shouldFail && err == nil {
				t.Fatal("Expected to fail but err is nil")
			}

			if diff := deep.Equal(c, tt.config); diff != nil {
				t.Fatal(diff)
			}
		})
	}

}
