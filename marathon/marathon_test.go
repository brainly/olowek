package marathon

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-test/deep"
)

func TestCreateMarathonClient(t *testing.T) {
	c, err := NewMarathonClient("http://127.0.0.1:8080")

	if err != nil {
		t.Fatalf("Unexpected error: '%s'", err)
	}

	if c == nil {
		t.Fatalf("Client should not be nil")
	}
}

func TestCreateMarathonClientFailOnInvalidUrl(t *testing.T) {
	c, err := NewMarathonClient("this_is_invalid")

	if err == nil {
		t.Fatalf("Expected NewMarathonClient to return error")
	}

	if c != nil {
		t.Fatalf("Client instance should be nil on error")
	}

}

func TestGetApplications(t *testing.T) {
	getAppsCases := map[string]struct {
		fixture string
		scope   string
		apps    []Application
	}{
		"No apps from Marathon should give empty Application array": {
			fixture: "./fixtures/no_apps.json",
			apps:    []Application{},
		},
		"Apps with mixed scope and healthy/unhealthy/invalid tasks should produce valid list of tasks": {
			fixture: "./fixtures/apps.json",
			scope:   "",
			apps: []Application{
				Application{
					Name: "foo",
					ID:   "/production/foo",
					Labels: map[string]string{
						"scope": "public",
					},
					Env: map[string]string{
						"SERVICE_ENVIRONMENT": "production",
					},
					Tasks: []Task{
						Task{
							ID:    "production_foo.task1",
							Host:  "127.0.0.1",
							Ports: []int{5411},
						},
					},
				},
				Application{
					Name: "bar",
					ID:   "/bar",
					Labels: map[string]string{
						"scope": "internal",
					},
					Env: map[string]string{},
					Tasks: []Task{
						Task{
							ID:    "bar.task3",
							Host:  "127.0.0.3",
							Ports: []int{5413},
						},
					},
				},
				Application{
					Name:   "baz",
					ID:     "/foo/bar/baz",
					Labels: map[string]string{},
					Env:    map[string]string{},
					Tasks:  []Task{},
				},
			},
		},
		"Should have only apps with scope=internal": {
			fixture: "./fixtures/apps.json",
			scope:   "internal",
			apps: []Application{
				Application{
					Name: "bar",
					ID:   "/bar",
					Labels: map[string]string{
						"scope": "internal",
					},
					Env: map[string]string{},
					Tasks: []Task{
						Task{
							ID:    "bar.task3",
							Host:  "127.0.0.3",
							Ports: []int{5413},
						},
					},
				},
			},
		},
	}
	for name, tt := range getAppsCases {
		t.Run(name, func(t *testing.T) {
			buf, err := ioutil.ReadFile(tt.fixture)
			if err != nil {
				t.Fatalf("Error reading fixture file: '%s'", err)
			}
			response := string(buf)

			server := newFakeMarathonAppsServer(response)
			defer server.Close()

			c, err := NewMarathonClient(server.URL)
			if err != nil {
				t.Fatalf("Unexpected error: '%s'", err)
			}

			apps, err := c.GetApplications(tt.scope)
			if err != nil {
				t.Fatalf("Unexpected error: '%s'", err)
			}

			if diff := deep.Equal(tt.apps, apps); diff != nil {
				t.Fatal(diff)
			}
		})
	}

}

func newFakeMarathonAppsServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/v2/apps" {
			fmt.Fprintln(w, response)
		}
	}))
}
