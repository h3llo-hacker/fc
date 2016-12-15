package handler

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func ListServices(endpoint string) ([]swarm.Service, error) {
	log.Info(fmt.Sprintf("Get [ %s ] Services", endpoint))
	host := "tcp://" + endpoint
	version := "v1.24"
	UA := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(host, version, nil, UA)
	if err != nil {
		log.Error(err)
	}
	ctx := context.Background()

	// Get Nodes
	services, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err == nil {
		return services, nil
	} else {
		return []swarm.Service{}, err
	}
}
