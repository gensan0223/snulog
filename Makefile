run:
	go run main.go $(ARGS)

run-fetch:
	make run ARGS="fetch"

build:
	go build -o snulog main.go

run-server:
	go run mockserver/server.go
