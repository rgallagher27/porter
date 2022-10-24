.PHONY: build

run-porter:
	docker-compose up

test-porter:
	go test -v ./...

build:
	CGO_ENABLED="0" GOOS="linux" GOARCH="amd64" go build -o ./porter ./.