package stats

import (
	"encoding/json"
	"testing"
	"time"
)

type statsJSON struct {
	ReloadFailed   int `json:"reloads_failed"`
	ReloadOk       int `json:"reloads_ok"`
	RenderFailed   int `json:"renders_failed"`
	MarathonFailed int `json:"marathon_failed"`

	LastReload    string `json:"last_reload"`
	LastReloadAgo int    `json:"last_reload_ago"`
	LastRender    string `json:"last_render"`
	LastRenderAgo int    `json:"last_render_ago"`
	LastEvent     string `json:"last_event"`
	LastEventAgo  int    `json:"last_event_ago"`
}

func TestStatsMarshaling(t *testing.T) {
	s := NewStats()

	// Add some data
	s.ReloadFailed()
	s.RenderFailed()

	s.MarathonFailed()
	s.MarathonFailed()
	s.MarathonFailed()

	s.NginxReloaded()
	s.NginxReloaded()

	s.UpdateLastEvent()
	s.UpdateLastRender()

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Unexpected error when marshaling to JSON: '%s'", err)
	}

	newData := statsJSON{}
	err = json.Unmarshal(data, &newData)
	if err != nil {
		t.Fatalf("Unexpected error when unmarshaling from JSON: '%s'", err)
	}

	if newData.ReloadFailed != 1 {
		t.Fatalf("Expected to have 1 reload failed stat")
	}

	if newData.ReloadOk != 2 {
		t.Fatalf("Expected to have 2 reload ok stats")
	}

	if newData.MarathonFailed != 3 {
		t.Fatalf("Expected to have 3 marathon failed stats")
	}

	if newData.RenderFailed != 1 {
		t.Fatalf("Expected to have 1 render failed stat")
	}

	var emptyTime time.Time
	emptyTimeFormat := emptyTime.Format(time.RFC3339)

	if newData.LastReload == emptyTimeFormat {
		t.Fatalf("LastReload should be set but is: '%s'", emptyTimeFormat)
	}

	if newData.LastRender == emptyTimeFormat {
		t.Fatalf("LastRender should be set but is: '%s'", emptyTimeFormat)
	}

	if newData.LastEvent == emptyTimeFormat {
		t.Fatalf("LastEvent should be set but is: '%s'", emptyTimeFormat)
	}
}
