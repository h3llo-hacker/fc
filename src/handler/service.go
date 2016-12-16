package handler

import (
	"config"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func ListServices(endpoint string) ([]string, error) {
	log.Info(fmt.Sprintf("Get [ %s ] Services", endpoint))
	host := "tcp://" + endpoint
	version := "v1.24"
	UA := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(host, version, nil, UA)
	if err != nil {
		log.Error(err)
	}
	ctx := context.Background()

	// Get Services
	services, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err == nil {
		s := make([]string, 0)
		for _, service := range services {
			s = append(s, service.Spec.Name)
		}
		return s, nil
	} else {
		return make([]string, 1), err
	}
}

func InspectService(serviceID string) (swarm.Service, error) {
	log.Info(fmt.Sprintf("Get Service [ %s ]", serviceID))
	for _, endpoint := range config.Conf.Endpoints {
		host := "tcp://" + endpoint
		version := "v1.24"
		UA := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		cli, err := client.NewClient(host, version, nil, UA)
		if err != nil {
			log.Error(err)
		}
		ctx := context.Background()

		// Get Service
		// filters := filters.NewArgs()
		// filters.Add("service", serviceID)
		service, _, err := cli.ServiceInspectWithRaw(ctx, serviceID)
		if err == nil {
			return service, nil
		}
	}
	return swarm.Service{}, nil
}

func InspectServiceTasks(serviceID string) (swarm.Task, error) {
	for _, endpoint := range config.Conf.Endpoints {
		host := "tcp://" + endpoint
		version := "v1.24"
		UA := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		cli, err := client.NewClient(host, version, nil, UA)
		if err != nil {
			log.Error(err)
		}
		ctx := context.Background()
		// Get Service
		service, _, err := cli.ServiceInspectWithRaw(ctx, serviceID)
		if err == nil {
			filters := filters.NewArgs()
			filters.Add("service", service.ID)
			filters.Add("desired-state", "running")
			tasks, err := cli.TaskList(ctx, types.TaskListOptions{Filters: filters})
			if err != nil {
				log.Error(err)
			}
			// Normally, there is only one task, but in global mode, there are lots of tasks.
			var v_max uint64
			var task swarm.Task
			for _, t := range tasks {
				if t.Version.Index > v_max {
					v_max = t.Version.Index
					task = t
				}
			}
			return task, nil
		}
	}
	return swarm.Task{}, nil
}