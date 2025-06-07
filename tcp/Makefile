.PHONY: default check test build image

IMAGE_NAME := teebow1e/whoami-tcp

default: check test build

build:
	CGO_ENABLED=0 go build -a --trimpath --installsuffix cgo --ldflags="-s" -o whoami

check:
	golangci-lint run

image:
	docker build -t $(IMAGE_NAME) .
