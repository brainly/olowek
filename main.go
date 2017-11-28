package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/marathon"
	"github.com/brainly/olowek/utils"
	"github.com/brainly/olowek/worker"
)

const (
	DefaultConfigPath = "/etc/olowek/olowek.json"
)

func main() {
	cfgFlag := flag.String("c", DefaultConfigPath, fmt.Sprintf("Path to configuration file [default: %s]", DefaultConfigPath))
	flag.Parse()

	cfg, err := config.NewConfigFromFile(*cfgFlag)
	panicOnError(err)

	cfg.NginxReloadFunc = utils.NginxReload

	client := setupMarathon(cfg.Marathon)

	worker := worker.Worker{
		Trigger: make(chan bool, 2),
		Action:  worker.NewNginxReloaderWorker(client, cfg),
	}
	worker.Work()

	client.ConnectToEventStream(worker.Trigger)
}

func setupMarathon(marathonURL string) marathon.Marathon {
	client, err := marathon.NewMarathonClient(marathonURL)
	panicOnError(err)

	return client
}

func panicOnError(err error) {
	if err != nil {
		log.Panicf("Panic: %v", err)
	}
}
