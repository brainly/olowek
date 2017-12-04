package worker

import (
	"reflect"
	"time"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/marathon"
	"github.com/brainly/olowek/utils"
	log "github.com/sirupsen/logrus"
)

func NewNginxReloaderWorker(client marathon.Marathon, cfg *config.Config) func() {
	return func() {
		start := time.Now()
		cfg.Lock()
		defer cfg.Unlock()

		log.Info("Updating nginx config")
		apps, err := client.GetApplications(cfg.Scope)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Error getting applications from Marathon")
			return
		}

		if reflect.DeepEqual(cfg.Apps, apps) {
			log.Info("No changes in configuration")
			return
		}

		cfg.Apps = apps

		err = utils.RenderTemplate(cfg.NginxTemplate, cfg.NginxConfig, cfg)
		if err != nil {
			log.WithFields(log.Fields{
				"src":  cfg.NginxTemplate,
				"dest": cfg.NginxConfig,
				"err":  err,
			}).Error("Error generating template")
			return
		}

		log.Info("Reloading nginx")
		err = cfg.NginxReloadFunc(cfg.NginxCmd)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("Error reloading nginx")
			return
		}

		elapsed := time.Since(start)
		log.WithFields(log.Fields{"took": elapsed}).Info("Generated new nginx configuration")
	}
}
