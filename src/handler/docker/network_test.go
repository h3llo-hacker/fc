package docker

import (
	"testing"
)

var endpoint = "127.0.0.1:2374"

func TestListNetworks(t *testing.T) {
	n, err := ListNetworks(endpoint)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(n)
	}
}

func TestCreateNetwork(t *testing.T) {
	networkID := "test-net"
	n, err := CreateNetwork(endpoint, networkID)
	if err == nil && n == networkID {
		t.Log("Create Network OK")
	} else {
		t.Error(err)
	}
}

func TestInspectNetwork(t *testing.T) {
	networkID := "test-net"
	N, err := InspectNetwork(endpoint, networkID)
	if err == nil && N.Name == networkID {
		t.Log(N.ID)
		t.Log("Inspect Network OK")
	} else {
		t.Error(err)
	}
}

func TestRemoveNetwork(t *testing.T) {
	networkID := "test-net"
	err := RemoveNetwork(endpoint, networkID)
	if err == nil {
		t.Log("Remove Network OK")
	} else {
		t.Error(err)
	}
}
