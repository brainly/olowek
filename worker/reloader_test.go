package worker

import (
	"testing"

	marathon "github.com/gambol99/go-marathon"
)

func TestIsTaskHealthy(t *testing.T) {
	isTaskHealthyCases := map[string]struct {
		task   marathon.Task
		result bool
	}{
		"Missing ports should be unhealthy": {
			task: marathon.Task{
				Host:  "10.0.0.1",
				Ports: []int{},
			},
			result: false,
		},
		"Missing host should be unhealthy": {
			task: marathon.Task{
				Ports: []int{5411},
			},
			result: false,
		},
		"Missing healthchecks should be healthy": {
			task: marathon.Task{
				Host:               "10.0.0.1",
				Ports:              []int{5411},
				HealthCheckResults: nil,
			},
			result: true,
		},
		"All healthchecks alive should be healthy": {
			task: marathon.Task{
				Host:  "10.0.0.1",
				Ports: []int{5411},
				HealthCheckResults: []*marathon.HealthCheckResult{
					{
						Alive: true,
					},
					{
						Alive: true,
					},
				},
			},
			result: true,
		},
		"One failed health check should be unhealthy": {
			task: marathon.Task{
				Host:  "10.0.0.1",
				Ports: []int{5411},
				HealthCheckResults: []*marathon.HealthCheckResult{
					{
						Alive: true,
					},
					{
						Alive: false,
					},
				},
			},
			result: false,
		},
	}

	for name, tt := range isTaskHealthyCases {
		t.Run(name, func(t *testing.T) {
			result := isTaskHealthy(tt.task)

			if tt.result != result {
				t.Fatalf("Unexpected result: expected '%v', got '%v'", tt.result, result)
			}
		})
	}
}

func TestIsInScope(t *testing.T) {
	inScopeCases := map[string]struct {
		app    marathon.Application
		scope  string
		result bool
	}{
		"Empty scope should pass": {
			app:    marathon.Application{},
			scope:  "",
			result: true,
		},
		"Empty lables should filter out": {
			app:    marathon.Application{},
			scope:  "public",
			result: false,
		},
		"No scope labe should filter out": {
			app: marathon.Application{
				Labels: &map[string]string{
					"foo": "bar",
				},
			},
			scope:  "public",
			result: false,
		},
		"Different scope should filter out": {
			app: marathon.Application{
				Labels: &map[string]string{
					"scope": "internal",
				},
			},
			scope:  "public",
			result: false,
		},
		"Matching scope should pass (1)": {
			app: marathon.Application{
				Labels: &map[string]string{
					"scope": "public",
				},
			},
			scope:  "public",
			result: true,
		},
		"Matching scope should pass (2)": {
			app: marathon.Application{
				Labels: &map[string]string{
					"scope": "internal",
				},
			},
			scope:  "internal",
			result: true,
		},
	}

	for name, tt := range inScopeCases {
		t.Run(name, func(t *testing.T) {
			result := isInScope(tt.app, tt.scope)

			if tt.result != result {
				t.Fatalf("Unexpected result: expected '%v', got '%v'", tt.result, result)
			}
		})
	}
}
