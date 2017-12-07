package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/brainly/olowek/api"
	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/marathon"
	"github.com/brainly/olowek/stats"
	"github.com/brainly/olowek/utils"
	"github.com/brainly/olowek/worker"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultConfigPath = "/etc/olowek/olowek.json"
)

var (
	VERSION = "master"
)

var opts struct {
	Config  string `short:"c" long:"config" description:"Path to configuration file." default:"/etc/olowek/olowek.json"`
	Debug   bool   `short:"d" long:"debug" description:"Enable debug logging"`
	Version func() `short:"v" long:"version" description:"Print version and exit"`
}

func versionFunc() {
	fmt.Printf("olowek\nVersion: %s\n", VERSION)
	os.Exit(0)
}

func main() {
	opts.Version = versionFunc

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

	s := stats.NewStats()

	client, err := marathon.NewMarathonClient(cfg.Marathon)
	if err != nil {
		log.WithFields(log.Fields{
			"url": cfg.Marathon,
			"err": err,
		}).Fatal("Error creating Marathon client")
	}

	go setupHttpServer(cfg, s)
	log.WithFields(log.Fields{
		"addr": cfg.BindAddress,
	}).Info("Started http server")

	worker := worker.Worker{
		Trigger: make(chan bool, 2),
		Action:  worker.NewNginxReloaderWorker(client, cfg, s),
	}
	worker.Work()

	client.ConnectToEventStream(worker.Trigger)
}

func setupHttpServer(cfg *config.Config, s stats.Stats) {
	r := mux.NewRouter()
	r.HandleFunc("/v1/stats", api.StatsHandler(s))

	err := http.ListenAndServe(cfg.BindAddress, r)
	if err != nil {
		log.WithFields(log.Fields{
			"addr": cfg.BindAddress,
			"err":  err,
		}).Fatal("Failed to start http server")
	}
}
