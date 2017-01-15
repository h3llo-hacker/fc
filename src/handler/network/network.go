package handler

import (
	// "config"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	// "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	// "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func ListNetworks(endpoint string) ([]string, error) {
	log.Info(fmt.Sprintf("Get [ %s ] Networks", endpoint))
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return []string{}, err
	} else {
		defer cli.Close()
	}

	// Get Networks
	ctx := context.Background()
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err == nil {
		n := make([]string, 0)
		for _, network := range networks {
			n = append(n, network.Name)
		}
		log.Debug(n)
		return n, nil
	} else {
		return make([]string, 1), err
	}
}

func CreateNetwork(endpoint, networkID string) (string, error) {
	log.Info(fmt.Sprintf("Create Network [ %s ]", networkID))
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return "", err
	} else {
		defer cli.Close()
	}

	// Check
	if HasNetwork(cli, networkID) {
		return "", errors.New("Network Already exists.")
	}

	// Create Networks
	ctx := context.Background()
	Cnetwork, err := cli.NetworkCreate(ctx, networkID, types.NetworkCreate{
		Driver: "overlay",
		IPAM: &network.IPAM{
			Driver: "default",
		},
	})
	if err == nil {
		n, err := cli.NetworkInspect(ctx, Cnetwork.ID)
		if err != nil {
			return "", err
		}
		return n.Name, nil
	} else {
		return "", err
	}

}

func RemoveNetwork(endpoint, networkID string) error {
	log.Info(fmt.Sprintf("Remove Network [ %s ]", networkID))
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return err
	} else {
		defer cli.Close()
	}

	// Check
	if !HasNetwork(cli, networkID) {
		return errors.New("Network 404.")
	}

	// Remove Networks
	ctx := context.Background()
	err = cli.NetworkRemove(ctx, networkID)
	if err == nil {
		return nil
	} else {
		return err
	}

}

func InspectNetwork(endpoint, networkID string) (*types.NetworkResource, error) {
	log.Info(fmt.Sprintf("Inspect Network [ %s ]", networkID))
	cli, err := DockerCli(endpoint)
	if err != nil {
		log.Error(err)
		return &types.NetworkResource{}, err
	} else {
		defer cli.Close()
	}

	// Check
	if !HasNetwork(cli, networkID) {
		return &types.NetworkResource{}, errors.New("Network 404.")
	}

	// Inspect Networks
	ctx := context.Background()
	Inetwork, err := cli.NetworkInspect(ctx, networkID)
	if err == nil {
		return &Inetwork, nil
	} else {
		return &types.NetworkResource{}, err
	}
}

func HasNetwork(cli *client.Client, networkName string) bool {
	// Get Networks
	ctx := context.Background()
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err == nil {
		for _, network := range networks {
			if network.Name == networkName {
				return true
			}
		}
		return false
	} else {
		log.Error(err)
		return false
	}
}
