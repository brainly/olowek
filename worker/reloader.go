package worker

import (
	"log"
	"reflect"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/marathon"
	"github.com/brainly/olowek/utils"
)

func NewNginxReloaderWorker(client marathon.Marathon, cfg *config.Config) func() {
	return func() {
		cfg.Lock()
		defer cfg.Unlock()

		log.Printf("Updating nginx config")
		apps, err := client.GetApplications(cfg.Scope)
		if err != nil {
			log.Printf("Error getting applications: '%s'", err)
			return
		}

		if reflect.DeepEqual(cfg.Apps, apps) {
			log.Printf("No changes in configuration")
			return
		}

		cfg.Apps = apps

		err = utils.RenderTemplate(cfg.NginxTemplate, cfg.NginxConfig, cfg)
		if err != nil {
			log.Printf("Error generating template: '%s'", err)
			return
		}

		log.Printf("Reloading nginx")
		err = cfg.NginxReloadFunc(cfg.NginxCmd)
		if err != nil {
			log.Printf("Error reloading nginx: '%s'", err)
			return
		}
	}
}
