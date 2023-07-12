build:
	go build -o bin/larkdown main.go

run:
	go run main.go

test:
	go test -v ./...