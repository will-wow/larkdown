default: test

build:
	go build -o bin/larkdown main.go

run:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	goimports -w -local github.com/will-wow/larkdown .
