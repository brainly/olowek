package worker

import (
	"log"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/marathon"
	"github.com/brainly/olowek/utils"
)

func NewNginxReloaderWorker(client marathon.Marathon, cfg *config.Config) func() {
	return func() {
		cfg.Lock()
		defer cfg.Unlock()

		log.Printf("Updating nginx config")
		err := generateNginxConfig(client, cfg)
		if err != nil {
			log.Printf("Error generating template: '%s'", err)
			return
		}

		err = cfg.NginxReloadFunc(cfg.NginxCmd)
		if err != nil {
			log.Printf("Error reloading nginx: '%s'", err)
			return
		}
	}
}

func generateNginxConfig(client marathon.Marathon, cfg *config.Config) error {
	apps, err := client.GetApplications(cfg.Scope)
	if err != nil {
		return err
	}

	cfg.Apps = apps

	err = utils.RenderTemplate(cfg.NginxTemplate, cfg.NginxConfig, cfg)
	if err != nil {
		return err
	}

	return nil
}
