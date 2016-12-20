package handler

import (
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

func TestInspectService(t *testing.T) {
	serviceID := "nginx"
	service, err := InspectService(serviceID)
	if err == nil {
		if service.Spec.Name == serviceID {
			t.Log("OK")
		} else {
			t.Error(service.ID)
		}
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
