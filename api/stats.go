package api

import (
	"encoding/json"
	"net/http"

	"github.com/brainly/olowek/stats"
)

func StatsHandler(s stats.Stats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := json.Marshal(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
