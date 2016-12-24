package handler

import (
	"config"
	"errors"
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
	} else {
		defer cli.Close()
	}
	// Get Services
	ctx := context.Background()
	services, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	defer cli.Close()
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
	var S swarm.Service
	var E error
	for _, endpoint := range config.Conf.Endpoints {
		log.Info(fmt.Sprintf("Get Service [ %s ]", serviceID))
		cli, err := DockerCli(endpoint)
		if err != nil {
			log.Error(err)
			return swarm.Service{}, err
		} else {
			defer cli.Close()
		}

		// Get Service
		ctx := context.Background()
		service, _, err := cli.ServiceInspectWithRaw(ctx, serviceID)
		if err == nil {
			S = service
			E = nil
			break
		} else {
			S = swarm.Service{}
			E = err
		}
	}
	return S, E
}

func InspectServiceTasks(serviceID string) (swarm.Task, error) {
	for _, endpoint := range config.Conf.Endpoints {
		log.Info(endpoint)
		cli, err := DockerCli(endpoint)
		if err != nil {
			log.Error(err)
			return swarm.Task{}, err
		}

		// Get Service
		ctx := context.Background()
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
			defer cli.Close()
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

	// Check service
	if HasService(cli, serviceName) == true {
		e := errors.New(fmt.Sprintf("Service [ %s ] Already Exists.", serviceName))
		return e
	}

	// Create Service
	ctx := context.Background()
	response, err := cli.ServiceCreate(ctx, *service, types.ServiceCreateOptions{})
	defer cli.Close()
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
	if HasService(cli, serviceID) == false {
		e := errors.New(fmt.Sprintf("Service [ %s ] Not Found", serviceID))
		return e
	}
	err = cli.ServiceRemove(ctx, serviceID)
	defer cli.Close()
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

func HasService(cli *client.Client, serviceID string) bool {
	ctx := context.Background()
	_, _, err := cli.ServiceInspectWithRaw(ctx, serviceID)
	defer cli.Close()
	if err != nil {
		return false
	}
	return true
}
