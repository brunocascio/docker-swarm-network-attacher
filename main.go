package main

import (
	"context"
	"log"

	"brunocascio/docker-swarm-network-attacher/lib"

	"github.com/docker/docker/client"
)

func main() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	log.Printf("Connected to docker endpoint %s (version %s)", cli.DaemonHost(), cli.ClientVersion())

	listeners, err := lib.GetListeners(cli, ctx)

	if err != nil {
		log.Fatalln(err)
	}

	targets, err := lib.GetTargets(cli, ctx)

	if err != nil {
		log.Fatalln(err)
	}

	for _, listener := range listeners {

		listener_ignore_networks := lib.GetListenerIgnoreNetworks(cli, ctx, listener)
		listener_targets_networks := lib.GetListenerTargetsNetworks(cli, ctx, listener, targets)

		networks_to_update := append(listener_ignore_networks, listener_targets_networks...)

		lib.UpdateListenerNetworks(cli, ctx, listener, networks_to_update)
	}
}
