package main

import (
	"os"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/marathon"
	"github.com/brainly/olowek/utils"
	"github.com/brainly/olowek/worker"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultConfigPath = "/etc/olowek/olowek.json"
)

var opts struct {
	Config string `short:"c" long:"config" description:"Path to configuration file." default:"/etc/olowek/olowek.json"`
	Debug  bool   `short:"d" long:"debug" description:"Enable debug logging"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}

	cfg, err := config.NewConfigFromFile(opts.Config)
	if err != nil {
		log.WithFields(log.Fields{
			"config": opts.Config,
			"err":    err,
		}).Fatal("Error reading configuration from file")
	}
	cfg.NginxReloadFunc = utils.NginxReload

	client, err := marathon.NewMarathonClient(cfg.Marathon)
	if err != nil {
		log.WithFields(log.Fields{
			"url": cfg.Marathon,
			"err": err,
		}).Fatal("Error creating Marathon client")
	}

	worker := worker.Worker{
		Trigger: make(chan bool, 2),
		Action:  worker.NewNginxReloaderWorker(client, cfg),
	}
	worker.Work()

	client.ConnectToEventStream(worker.Trigger)
}
