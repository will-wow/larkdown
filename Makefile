default: test

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	goimports -w -local github.com/will-wow/larkdown .

ready: fmt lint test