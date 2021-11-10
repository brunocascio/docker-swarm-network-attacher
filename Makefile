dev:
	docker build --target base -t dsna:local .
	docker run --rm -ti -v "${PWD}:/app" -w /app -v /var/run/docker.sock:/var/run/docker.sock dsna:local

build:
	@docker build -t dsna .

examples-simple:
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/subscriber.yml subscriber
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/services.yml services-1
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/services.yml services-2
	@docker stack deploy --prune --resolve-image=changed -c examples/simple/services.yml services-3

examples-clean:
	@docker stack rm subscriber services-{1,2,3}