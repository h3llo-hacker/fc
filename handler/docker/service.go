package docker

import (
	"errors"
	"fmt"

	"github.com/h3llo-hacker/fc/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func ListServices(endpoint string) ([]string, error) {
	log.Infof("Get [ %s ] Services", endpoint)
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
	if err != nil {
		return make([]string, 1), err
	}
	s := make([]string, len(services))
	for _, service := range services {
		s = append(s, service.Spec.Name)
	}
	return s, nil
}

func InspectService(serviceID string) (swarm.Service, error) {
	var (
		S swarm.Service
		E error
	)
	endpoint := config.Conf.Endpoint

	log.Debugf(fmt.Sprintf("Get Service [ %s ]", serviceID))
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
	if err != nil {
		S = swarm.Service{}
		E = err
	}
	S = service
	E = nil
	return S, E
}

func InspectServiceTasks(serviceID string) (swarm.Task, error) {
	endpoint := config.Conf.Endpoint
	log.Info(endpoint)
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return swarm.Task{}, err
	}

	// Get Service
	ctx := context.Background()
	service, _, err := cli.ServiceInspectWithRaw(ctx, serviceID)
	if err != nil {
		return swarm.Task{}, nil
	}
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
		return errors.New(fmt.Sprintf("Service [ %s ] Already Exists.", serviceName))
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
