package stats

import (
	"encoding/json"
	"sync"
	"time"
)

type Stats interface {
	// Update counters
	ReloadFailed()
	RenderFailed()
	MarathonFailed()

	NginxReloaded()

	// Succesful actions dates
	UpdateLastEvent()
	UpdateLastRender()

	MarshalJSON() ([]byte, error)
}

type stats struct {
	sync.RWMutex

	failedReloads  int
	failedRenders  int
	failedMarathon int
	nginxReloads   int

	lastRender time.Time
	lastReload time.Time
	lastEvent  time.Time
}

func NewStats() Stats {
	return &stats{}
}

func (s *stats) MarshalJSON() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()
	data := struct {
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
	}{
		ReloadFailed: s.failedReloads,
		ReloadOk:     s.nginxReloads,

		RenderFailed:   s.failedRenders,
		MarathonFailed: s.failedMarathon,

		LastReload:    s.lastReload.Format(time.RFC3339),
		LastReloadAgo: int(time.Since(s.lastReload).Seconds()),
		LastRender:    s.lastRender.Format(time.RFC3339),
		LastRenderAgo: int(time.Since(s.lastRender).Seconds()),
		LastEvent:     s.lastEvent.Format(time.RFC3339),
		LastEventAgo:  int(time.Since(s.lastEvent).Seconds()),
	}

	return json.Marshal(data)
}

func (s *stats) ReloadFailed() {
	s.Lock()
	defer s.Unlock()
	s.failedReloads += 1
}

func (s *stats) RenderFailed() {
	s.Lock()
	defer s.Unlock()
	s.failedRenders += 1
}

func (s *stats) MarathonFailed() {
	s.Lock()
	defer s.Unlock()
	s.failedMarathon += 1
}

func (s *stats) NginxReloaded() {
	s.Lock()
	defer s.Unlock()
	s.nginxReloads += 1
	s.lastReload = time.Now().UTC()
}

func (s *stats) UpdateLastEvent() {
	s.Lock()
	defer s.Unlock()
	s.lastEvent = time.Now().UTC()
}

func (s *stats) UpdateLastRender() {
	s.Lock()
	defer s.Unlock()
	s.lastRender = time.Now().UTC()
}
