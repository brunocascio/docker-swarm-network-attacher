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

type Connection struct {
	connection_name string
	network_name    string
}

func Start(cli *client.Client, ctx context.Context) {

	connections, err := GetConnections(cli, ctx)

	if err != nil {
		log.Fatalln(err)
		return
	}

	hydrated_connections := HydrateConnections(cli, ctx, connections)

	for connection, networks := range hydrated_connections {
		err := UpdateConnectionNetworks(cli, ctx, connection, networks)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func UpdateConnectionNetworks(cli *client.Client, ctx context.Context, connection string, networks []string) error {
	var changed = false
	service, _, err := cli.ServiceInspectWithRaw(ctx, connection, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}
	current_svc_networks := funk.Get(service.Spec.TaskTemplate.Networks, "Target")
	funk.ForEach(networks, func(net string) {
		if !funk.Contains(current_svc_networks, net) {
			changed = true
			service.Spec.TaskTemplate.Networks = append(service.Spec.TaskTemplate.Networks, swarm.NetworkAttachmentConfig{Target: net})
		}
	})
	if changed {
		log.Printf("Updating service %s with networks %s", service.Spec.Name, service.Spec.TaskTemplate.Networks)
	}
	r, err := cli.ServiceUpdate(ctx, service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	_ = r
	return err
}

func HydrateConnections(cli *client.Client, ctx context.Context, connections []Connection) map[string][]string {
	var dict = make(map[string][]string)
	for _, connection := range connections {
		connection_obj, _, err := cli.ServiceInspectWithRaw(ctx, connection.connection_name, types.ServiceInspectOptions{})
		if err != nil {
			log.Fatalln(err)
			continue
		}
		network_obj, err := cli.NetworkInspect(ctx, connection.network_name, types.NetworkInspectOptions{})
		if err != nil {
			log.Fatalln(err)
			continue
		}
		dict[connection_obj.ID] = append(dict[connection_obj.ID], network_obj.ID)
	}
	return dict
}

func GetConnections(cli *client.Client, ctx context.Context) ([]Connection, error) {
	services, err := cli.ServiceList(ctx, types.ServiceListOptions{Filters: filters.NewArgs(filters.KeyValuePair{
		Key:   "label",
		Value: "dsna.connections",
	})})

	if err != nil {
		return nil, err
	}

	var connections []Connection

	for _, service := range services {
		connection_names := strings.Split(service.Spec.Labels["dsna.connections"], ",")
		for _, connection_name := range connection_names {
			connection_name := strings.TrimSpace(connection_name)
			network_name := service.Spec.Labels["dsna.connections."+connection_name+".network"]
			connections = append(connections, Connection{
				connection_name: connection_name,
				network_name:    network_name,
			})
		}
	}

	return connections, nil
}
