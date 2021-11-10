# Docker swarm network attacher

## Description

`docker-swarm-network-attacher` aims to solve the problem of sharing a network between unrelated containers. 
With this service we can generate "point-to-point" networks and avoid this problem.

**How it works**

It basically fetch **containers** (not services) with label `"dsna.enabled=true"` and auto connect them to all the overlay networks available on the cluster.
Under the hood it performs a `docker network connect` on every container to every network.
## How to deploy it

### Deploy dsna service

**dsna.stack.yml**

```yaml
version: "3.9"

services:
  service:
    image: bcascio/docker-swarm-network-attacher:<VERSION>
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    deploy:
      replicas: 1
      resources:
        limits:
          memory: 20M
          cpus: '0.05'
        reservations:
          memory: 10M
          cpus: '0.05'
```

**Run the deployment**

```sh
docker stack deploy -c dsna.stack.yml dsna
```

## Use Cases

- Cluster Gateway/Ingress
- External services that require communication with other isolated services on different networks (e.g. prometheus)

## Examples

See `examples` folder to see the use cases on action

## How to contribute

**Run development**

```sh
make dev
```