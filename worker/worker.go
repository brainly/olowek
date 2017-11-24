package worker

import (
	"time"
)

type Worker struct {
	Trigger chan bool
	Action  func()
}

func (w *Worker) Work() {
	go func() {

		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				<-w.Trigger
				w.Action()
			}
		}
	}()
}
