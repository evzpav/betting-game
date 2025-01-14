include help.mk
include .env

.DEFAULT_GOAL := help

BUILD       = $(shell git rev-parse --short HEAD)
DATE        = $(shell date -uIseconds)
VERSION  	= $(shell git describe --always --tags)
NAME        = $(shell basename $(CURDIR))
IMAGE       = $(NAME):$(BUILD)

git-config:
	git config --replace-all core.hooksPath .githooks

check-env-%:
	@ if [ "${${*}}" = ""  ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

version: ##@other Check version.
	@echo $(VERSION)

test: ##@dev Run unit tests locally
	go test ./...

build-local: ##@dev Build binary locally
	-rm ./betting-game

	CGO_ENABLED=0 \
	GOOS=linux  \
	GOARCH=amd64  \
	go build -installsuffix cgo -o betting-game -ldflags \
	"-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" \
	./cmd/server/main.go

install: ##@dev Instal Vue frontend dependencies
	cd ./frontend && npm install

run-frontend: install ##@dev Run Vue frontend locally at port 8080
	cd ./frontend && npm run serve

build-frontend: ##@dev Build static files for Vue
	cd ./frontend && npm run build

run-local: build-local ##@dev Run locally.
	HOST=localhost \
	PORT=8787 \
	LOGGER_LEVEL=debug \
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

test-docker: ##@docker Run unit tests in docker
	make target TARGET=test

lint-docker: ##@docker Run linting in docker
	make target TARGET=lint

image: ##@docker Build docker image.
	make target TARGET=image

run-docker: image ##@docker Run docker container.
	docker run --rm \
		--name $(NAME) \
		--network=host \
		-e HOST=localhost \
		-e PORT=8888 \
		-e LOGGER_LEVEL=info \
		$(IMAGE)

