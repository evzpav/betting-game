include help.mk
include .env

.PHONY: run-local git-config version clean install lint env env-stop test cover build image tag push deploy run run-docker remove-docker image
.DEFAULT_GOAL := help

BUILD         			= $(shell git rev-parse --short HEAD)
DATE          			= $(shell date -uIseconds)
VERSION  	  			= $(shell git describe --always --tags)
NAME           			= $(shell basename $(CURDIR))
IMAGE          			= $(NAME):$(BUILD)

git-config:
	git config --replace-all core.hooksPath .githooks

check-env-%:
	@ if [ "${${*}}" = ""  ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

version: ##@other Check version.
	@echo $(VERSION)


env-stop: ##@environment Remove mysql container.
	-docker rm -vf mysql_$(NAME)

build-local: ##@dev Build binary locally
	-rm ./betting-game

	CGO_ENABLED=0 \
	GOOS=linux  \
	GOARCH=amd64  \
	go build -installsuffix cgo -o betting-game -ldflags \
	"-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" \
	./cmd/server/main.go

run-frontend: ## Run Vue frontend locally at port 8080
	cd ./frontend && npm run serve

build-frontend: ## Build static files for Vue
	cd ./frontend && npm run build

run-local: build-local ##@dev Run locally.
	HOST=localhost \
	PORT=8787 \
	./betting-game

target: 
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=$(TARGET) \
		--file= .

build: ##@build Build image.
	make target TARGET=image

run-docker: build ##@docker Run docker container.
	docker run --rm \
		--name $(NAME) \
		--network=host \
		-e HOST=localhost \
		-e PORT=8080 \
		-e LOGGER_LEVEL=debug \
		$(IMAGE)
