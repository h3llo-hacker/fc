package docker

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/h3llo-hacker/fc/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/compose/convert"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Stack struct {
	// Name is the name of the stack
	Name string
	// Services is the number of the services
	Services int
}

func ListStacks(endpoint string) ([]*Stack, error) {
	log.Debugf("List [ %s ] Stacks", endpoint)
	cli, err := DockerCli(endpoint)

	// List Stacks
	ctx := context.Background()
	filter := filters.NewArgs()
	services, err := cli.ServiceList(
		ctx,
		types.ServiceListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Stack, 0)
	for _, service := range services {
		labels := service.Spec.Labels
		name, ok := labels[convert.LabelNamespace]
		if !ok {
			return nil,
				fmt.Errorf("cannot get label %s for service %s",
					convert.LabelNamespace, service.ID)
		}
		ztack, ok := m[name]
		if !ok {
			m[name] = &Stack{
				Name:     name,
				Services: 1,
			}
		} else {
			ztack.Services++
		}
	}
	var stacks []*Stack
	for _, stack := range m {
		stacks = append(stacks, stack)
	}
	return stacks, nil
}

// endpoint, stackName, coposeFile string
func DeployStack(endpoint, composeFile, stackName string) (string, error) {
	// docker -H <host> stack deploy -c test-docker-compose.yml nginx
	if !config.PathExist(composeFile) {
		return "",
			fmt.Errorf("compose file [" + composeFile + "] not found")
	}

	command := fmt.Sprintf("docker -H %s stack deploy -c %s %s",
		endpoint, composeFile, stackName)
	log.Debugf("exec command [%v]", command)
	cmd := exec.Command("docker", "-H", endpoint, "stack", "deploy", "-c", composeFile, stackName)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	logs := out.String()
	if err != nil {
		log.Errorf("Deploy stack [%v] error: [%v]",
			stackName, err)
		return logs, err
	}
	return logs, nil
}

func RemoveStack(endpoint, stackName string) error {
	// docker -H <host> stack rm nginx
	cmd := exec.Command("docker", "-H",
		endpoint, "stack", "rm", stackName)

	// var out bytes.Buffer
	// cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorf("Remove stack [%v] error: [%v]",
			stackName, err)
		return err
	}
	return nil
}

func PsStack(endpoint, namespace string) ([]swarm.Task, error) {
	filter := filters.NewArgs()
	filter.Add("label", convert.LabelNamespace+"="+namespace)
	filter.Add("desired-state", "running")
	client, err := DockerCli(endpoint)
	if err != nil {
		return nil, err
	}
	tasks, err := client.TaskList(context.Background(),
		types.TaskListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("Nothing found in stack: %s\n", namespace)
	}

	return tasks, nil
}
