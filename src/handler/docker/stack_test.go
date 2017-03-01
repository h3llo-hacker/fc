package docker

import (
	"testing"
	"time"
)

func Test_DeployStack(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	composeFile := "/home/mr/test/docker/test-docker-compose.yml"
	stackName := "nginx"
	t.Log(time.Now())
	deploylogs, err := DeployStack(endpoint, composeFile, stackName)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(deploylogs)
	t.Log(time.Now())
}

func Test_ListStacks(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	stacks, err := ListStacks(endpoint)
	if err != nil {
		t.Error(err)
	} else {
		for _, stack := range stacks {
			t.Log(stack.Name, stack.Services)
		}
	}
}

func Test_PsStack(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	stackName := "3c0ef85b-8201-49d9-5c06-6b486694a903"
	t.Log(time.Now())
	tasks, err := PsStack(endpoint, stackName)
	if err != nil {
		t.Error(err)
	}
	for _, task := range tasks {
		t.Log(task.ID)
	}
}

func Test_RemoveStack(t *testing.T) {
	endpoint := "127.0.0.1:2374"
	stackName := "nginx"
	t.Log(time.Now())
	err := RemoveStack(endpoint, stackName)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(time.Now())
}
