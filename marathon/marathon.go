package marathon

import (
	"net/url"
	"path"

	api "github.com/gambol99/go-marathon"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	Name   string
	ID     string
	Labels map[string]string
	Env    map[string]string
	Tasks  []Task
}

type Task struct {
	ID    string
	Host  string
	Ports []int
}

type Marathon interface {
	GetApplications(filterScope string) ([]Application, error)
	ConnectToEventStream(callback chan bool)
}

type client struct {
	client api.Marathon
}

func NewMarathonClient(marathonURL string) (Marathon, error) {
	config := api.NewDefaultConfig()
	config.URL = marathonURL
	config.EventsTransport = api.EventsTransportSSE

	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &client{
		client: c,
	}, nil
}

func (c *client) GetApplications(filterScope string) ([]Application, error) {
	filtered := []Application{}

	apps, err := c.client.Applications(url.Values{
		"embed": []string{"apps.tasks"},
	})

	if err != nil {
		return filtered, err
	}

	for _, a := range apps.Apps {
		if c.isInScope(a, filterScope) {
			var labels map[string]string
			if a.Labels != nil {
				labels = *a.Labels
			}

			var envs map[string]string
			if a.Env != nil {
				envs = *a.Env
			}

			filtered = append(filtered, Application{
				ID:     a.ID,
				Name:   path.Base(a.ID),
				Labels: labels,
				Env:    envs,
				Tasks:  c.getAppTasks(a),
			})
		}
	}

	return filtered, nil
}

func (c *client) ConnectToEventStream(callback chan bool) {
	// Register for events
	events, err := c.client.AddEventsListener(api.EventIDApplications)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("Failed to register for events")
	}

	log.Info("Doing full sync on start")
	callback <- true

	for {

		select {
		case event := <-events:
			select {
			case callback <- true:
				log.WithFields(log.Fields{
					"event": event,
				}).Debug("Received event. Pushing information to callback channel")
			default:
				log.Debug("Callback queue is full")
			}
		}
	}

	// Unsubscribe from Marathon events
	c.client.RemoveEventsListener(events)
}

func (c *client) getAppTasks(app api.Application) []Task {
	tasks := []Task{}

	if app.Tasks == nil {
		return tasks
	}

	for _, t := range app.Tasks {
		if !c.isTaskHealthy(*t) {
			continue
		}

		tasks = append(tasks, Task{
			ID:    t.ID,
			Host:  t.Host,
			Ports: t.Ports,
		})
	}

	return tasks
}

func (c *client) isTaskHealthy(task api.Task) bool {
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

func (c *client) isInScope(app api.Application, filterScope string) bool {
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
