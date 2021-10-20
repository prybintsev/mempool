build:
	go build -o bin/mempool cmd/main.go

test:
	go test -v ./...