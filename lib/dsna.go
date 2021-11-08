package lib

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/thoas/go-funk"
)

func GetListeners(cli *client.Client, ctx context.Context) ([]swarm.Service, error) {
	return cli.ServiceList(ctx, types.ServiceListOptions{Filters: filters.NewArgs(filters.KeyValuePair{
		Key:   "label",
		Value: "dsna.listener",
	})})
}

func GetTargets(cli *client.Client, ctx context.Context) ([]swarm.Service, error) {
	return cli.ServiceList(ctx, types.ServiceListOptions{Filters: filters.NewArgs(filters.KeyValuePair{
		Key:   "label",
		Value: "dsna.enable=true",
	})})
}

func UpdateListenerNetworks(cli *client.Client, ctx context.Context, listener swarm.Service, networks []string) (types.ServiceUpdateResponse, error) {
	networks_to_attach := funk.Map(networks, func(n string) swarm.NetworkAttachmentConfig {
		return swarm.NetworkAttachmentConfig{Target: n}
	})

	listener.Spec.TaskTemplate.Networks = networks_to_attach.([]swarm.NetworkAttachmentConfig)

	return cli.ServiceUpdate(ctx, listener.ID, listener.Version, listener.Spec, types.ServiceUpdateOptions{})
}

func GetListenerIgnoreNetworks(cli *client.Client, ctx context.Context, listener swarm.Service) []string {
	ignored_networks := strings.Split(strings.ReplaceAll(listener.Spec.Labels["dsna.listener.ignore-networks"], " ", ""), ",")
	var network_ids []string
	for _, n := range ignored_networks {
		network, err := cli.NetworkInspect(ctx, n, types.NetworkInspectOptions{})
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		network_ids = append(network_ids, network.ID)
	}
	return network_ids
}

func GetListenerTargetsNetworks(cli *client.Client, ctx context.Context, listener swarm.Service, targets []swarm.Service) []string {
	listener_name := listener.Spec.Name

	var network_ids []string

	for _, target := range targets {
		// remove whitespaces and convert sintrg into array
		target_listeners := strings.Split(strings.ReplaceAll(target.Spec.Labels["dsna.listeners"], " ", ""), ",")
		if !funk.Contains(target_listeners, listener_name) {
			// skip target since it's not related to the current listener
			continue
		}

		target_network := target.Spec.Labels["dsna.listeners."+listener_name+".network"]

		network, err := cli.NetworkInspect(ctx, target_network, types.NetworkInspectOptions{})

		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		network_ids = append(network_ids, network.ID)
	}

	return network_ids
}
