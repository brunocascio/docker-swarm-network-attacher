export CURRENT_DIR=${PWD}
export DOCKER_DEFAULT_PLATFORM=linux/amd64
export TARGETPLATFORM=${DOCKER_PLATFORM}

dev:
	@docker build --target base -t dnsa:local .
	@docker run --rm -ti -v "${CURRENT_DIR}:/app" -w /app -v /var/run/docker.sock:/var/run/docker.sock dnsa:local

build:
	@docker build -t dsna .

examples-simple:
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/a.yml simple-a
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/b.yml simple-b
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/c.yml simple-c
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/d.yml simple-d

examples-clean:
	@docker stack rm simple-{a,b,c,d}