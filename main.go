package main

import (
	"context"
	"log"
	"time"

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

	done := make(chan bool)
	go forever(cli, ctx)
	<-done // Block forever
}

func forever(cli *client.Client, ctx context.Context) {
	for {
		lib.Start(cli, ctx)
		time.Sleep(time.Second * 5)
	}
}
