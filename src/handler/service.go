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
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return []string{}, err
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
		cli, err := DockerCli(endpoint)
		if err != nil {
			log.Error(err)
			return swarm.Service{}, err
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
		cli, err := DockerCli(endpoint)
		if err != nil {
			log.Error(err)
			return swarm.Task{}, err
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

func CreateService(endpoint, serviceName, serviceImage string) error {
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return err
	}
	// Service Info
	service := &swarm.ServiceSpec{}
	service.Name = serviceName
	service.TaskTemplate.ContainerSpec.Image = serviceImage
	// Create Service
	ctx := context.Background()
	response, err := cli.ServiceCreate(ctx, *service, types.ServiceCreateOptions{})
	if err == nil {
		log.Info(response)
	} else {
		return err
	}
	return nil
}

func RemoveService(endpoint, serviceID string) error {
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return err
	}

	// Remove Service
	ctx := context.Background()
	err = cli.ServiceRemove(ctx, serviceID)
	if err != nil {
		return err
	}
	return nil
}

func DockerCli(endpoint string) (*client.Client, error) {
	host := "tcp://" + endpoint
	version := "v1.24"
	UA := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(host, version, nil, UA)
	if err != nil {
		log.Error(err)
		return &client.Client{}, err
	}
	return cli, nil
}
