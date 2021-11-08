DOCKER_HOST := tcp://172.17.10.103:2375

export DOCKER_HOST

dev:
	@./air main.go

build:
	@docker build -t dsna .