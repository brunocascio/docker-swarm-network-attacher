package lib

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func Start(cli *client.Client, ctx context.Context) {

	containers, err := GetSubscriberContainers(cli, ctx)

	if err != nil {
		log.Fatalln(err)
	}

	networks, err := GetNetworks(cli, ctx)

	if err != nil {
		log.Fatalln(err)
	}

	AttachContainersToNetworks(cli, ctx, containers, networks)
}

func AttachContainersToNetworks(cli *client.Client, ctx context.Context, containers []types.Container, networks []types.NetworkResource) {
	for _, container := range containers {
		for _, network := range networks {
			err := cli.NetworkConnect(ctx, network.ID, container.ID, nil)
			if err != nil {
				if strings.Contains(err.Error(), "already exists") {
					continue
				}
				log.Println(err, network)
			} else {
				log.Printf("attached %s network to %s container", network.Name, container.Names[0])
			}
		}
	}
}

func GetSubscriberContainers(cli *client.Client, ctx context.Context) ([]types.Container, error) {
	return cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(
		filters.KeyValuePair{
			Key:   "label",
			Value: "dsna.enable=true",
		},
		filters.KeyValuePair{
			Key:   "status",
			Value: "running",
		},
	),
	})
}

func GetNetworks(cli *client.Client, ctx context.Context) ([]types.NetworkResource, error) {
	return cli.NetworkList(ctx, types.NetworkListOptions{Filters: filters.NewArgs(filters.KeyValuePair{
		Key:   "driver",
		Value: "overlay",
	})})
}
