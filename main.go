package main

import (
	"log"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/utils"
	"github.com/brainly/olowek/worker"
	marathon "github.com/gambol99/go-marathon"
)

func main() {
	cfg, err := config.NewConfigFromFile("olowek.json")
	panicOnError(err)

	cfg.NginxReloadFunc = utils.NginxReload

	client := setupMarathon(cfg.Marathon)

	worker := worker.Worker{
		Trigger: make(chan bool, 2),
		Action:  worker.NewNginxReloaderWorker(client, cfg),
	}
	worker.Work()

	connectToEventStream(client, worker.Trigger)
}

func setupMarathon(marathonURL string) marathon.Marathon {
	config := marathon.NewDefaultConfig()
	config.URL = marathonURL
	config.EventsTransport = marathon.EventsTransportSSE

	client, err := marathon.NewClient(config)
	panicOnError(err)

	return client
}

func connectToEventStream(client marathon.Marathon, trigger chan bool) {
	// Register for events
	events, err := client.AddEventsListener(marathon.EventIDApplications)
	if err != nil {
		log.Fatalf("Failed to register for events, %s", err)
	}

	log.Printf("Force sync on start")
	trigger <- true

	for {

		select {
		case <-events:
			select {
			case trigger <- true:
			default:
				log.Printf("Callback queue is full")
			}
		}
	}

	// Unsubscribe from Marathon events
	client.RemoveEventsListener(events)
}

func panicOnError(err error) {
	if err != nil {
		log.Panicf("Panic: %v", err)
	}
}
