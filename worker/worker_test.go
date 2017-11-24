package worker

import (
	"testing"
)

func TestWorkerWillCallAction(t *testing.T) {
	testIterator := 0
	action := func() {
		testIterator += 1
	}

	w := Worker{
		Trigger: make(chan bool),
		Action:  action,
	}

	w.Work()
	w.Trigger <- true

	if testIterator != 1 {
		t.Fatalf("Expected testIterator to be '1', got '%d'", testIterator)
	}
}
