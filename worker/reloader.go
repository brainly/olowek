package worker

import (
	"log"
	"net/url"
	"path"

	"github.com/brainly/olowek/config"
	"github.com/brainly/olowek/utils"
	marathon "github.com/gambol99/go-marathon"
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
	apps, err := getAllApps(client, cfg.Scope)
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

func getAllApps(client marathon.Marathon, filterScope string) ([]config.App, error) {
	apps, err := client.Applications(url.Values{
		"embed": []string{"apps.tasks"},
	})

	if err != nil {
		return []config.App{}, err
	}

	var filtered []config.App
	for _, a := range apps.Apps {
		if isInScope(a, filterScope) {
			var labels map[string]string
			if a.Labels != nil {
				labels = *a.Labels
			}

			var envs map[string]string
			if a.Env != nil {
				envs = *a.Env
			}

			filtered = append(filtered, config.App{
				ID:     a.ID,
				Name:   path.Base(a.ID),
				Labels: labels,
				Env:    envs,
				Tasks:  getAppTasks(a),
			})
		}
	}

	return filtered, nil
}

func getAppTasks(app marathon.Application) []config.AppTask {
	var tasks []config.AppTask

	if app.Tasks == nil {
		return tasks
	}

	for _, t := range app.Tasks {
		if !isTaskHealthy(*t) {
			continue
		}

		tasks = append(tasks, config.AppTask{
			ID:           t.ID,
			Host:         t.Host,
			Ports:        t.Ports,
			ServicePorts: t.ServicePorts,
		})
	}

	return tasks
}

func isTaskHealthy(task marathon.Task) bool {
	// No ports exposed or missing Mesos Slave host address
	if len(task.Ports) == 0 || task.Host == "" {
		return false
	}

	// No health checks defind - nothing to check
	if task.HealthCheckResults == nil {
		return true
	}

	// Iterate all health checks to see if all are alive
	for _, health := range task.HealthCheckResults {
		if !health.Alive {
			return false
		}
	}

	return true
}

func isInScope(app marathon.Application, filterScope string) bool {
	if filterScope == "" {
		return true
	}

	if app.Labels != nil {
		labels := *app.Labels
		if scope, ok := labels["scope"]; ok && scope == filterScope {
			return true
		}
	}

	return false
}
