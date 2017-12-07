package worker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/marathon"
	"github.com/brainly/olowek/stats"
)

func TestNginxReloaderWorker(t *testing.T) {
	c, server := newFakeMarathonClient(t, "./fixtures/marathon.json")
	defer server.Close()

	tmpFile, err := ioutil.TempFile(".", ".services-test-")
	if err != nil {
		t.Fatalf("Unexpected error creating tmpfile: '%s'", err)
	}
	defer func() {
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Error closing tmpfile: '%s'", err)
		}

		os.Remove(tmpFile.Name())
	}()

	reloadFuncCalledTimes := 0
	cfg := &config.Config{
		Marathon:      server.URL,
		NginxConfig:   tmpFile.Name(),
		NginxTemplate: "./fixtures/services.tpl",
		NginxCmd:      "/bin/true",
		NginxReloadFunc: func(cmd string) error {
			reloadFuncCalledTimes += 1
			return nil
		},
	}
	s := stats.NewStats()

	reloader := NewNginxReloaderWorker(c, cfg, s)
	reloader()

	renderedTemplate, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error reading tmpfile: '%s'", err)
	}
	expectedConf, err := ioutil.ReadFile("./fixtures/services.conf")
	if err != nil {
		t.Fatalf("Unexpected error reading services.conf: '%s'", err)
	}

	if string(expectedConf) != string(renderedTemplate) {
		t.Fatalf("Rendered template is not as expected. Got:\n %s", string(renderedTemplate))
	}

	stat, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error while getting stat for tmpfile: '%s'", err)
	}
	modtime := stat.ModTime()

	// Sleep for 1s and try doing another worker call
	time.Sleep(time.Second)
	reloader()

	stat, err = os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error while getting stat for tmpfile: '%s'", err)
	}
	modtime_second_reload := stat.ModTime()

	if modtime != modtime_second_reload {
		t.Fatalf("File should not be modified since no configuration changes were made")
	}

	if reloadFuncCalledTimes != 1 {
		t.Fatalf("Reload func should be called only once since no configuration changes were made")
	}

}

func newFakeMarathonClient(t *testing.T, file string) (marathon.Marathon, *httptest.Server) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Error reading fixture file: '%s'", err)
	}

	server := newFakeMarathonAppsServer(string(buf))

	c, err := marathon.NewMarathonClient(server.URL)
	if err != nil {
		t.Fatalf("Unexpected error: '%s'", err)
	}

	return c, server
}

func newFakeMarathonAppsServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/v2/apps" {
			fmt.Fprintln(w, response)
		}
	}))
}
