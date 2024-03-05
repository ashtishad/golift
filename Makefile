run:
	export NUM_OF_SERVERS=5 \
	export STARTING_PORT=8000 \
	export LOAD_BALANCER_PORT=8080 \
	&& go run main.go
test:
	go test -v ./...
race:
	go run -race .
lint:
	golangci-lint run
