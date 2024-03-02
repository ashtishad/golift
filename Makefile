run:
	go run main.go
test:
	go test -v ./...
race:
	go run -race .
lint:
	golangci-lint run
