package docker

import (
	"errors"
	"testing"
)

func TestCreateService(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	serviceName := "nginx"
	serviceImage := "nginx"
	err := CreateService(endpoint, serviceName, serviceImage)
	if err == nil {
		t.Log("OK")
	} else {
		t.Error(err)
	}
}

func TestInspectServiceTasks(t *testing.T) {
	serviceID := "nginx"
	task, err := InspectServiceTasks(serviceID)
	if err == nil {
		t.Log("OK")
	} else {
		t.Error(err)
	}
	t.Log(task.Spec.ContainerSpec.Image)

}

func TestRemoveService(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	serviceID := "nginx"
	err := RemoveService(endpoint, serviceID)
	if err == nil {
		t.Log("OK")
	} else {
		t.Error(err)
	}
}

func TestHasService(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	cli, err := DockerCli(endpoint)
	defer cli.Close()
	if err != nil {
		t.Error(err)
	}
	serviceID := "nginxxxx"
	has := HasService(cli, serviceID)
	if has == true {
		t.Error(errors.New("service nginxxxx found?"))
	}
}

func TestDockerCli(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	cli, err := DockerCli(endpoint)
	// defer cli.Close()
	t.Log(cli.ClientVersion())
	if err != nil {
		t.Error(err)
	}
}
