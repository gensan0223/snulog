PROTO_SRC=proto/logs.proto

dev:
	docker compose up --build

run:
	go run main.go $(ARGS)

fetch:
	make run ARGS="fetch"

add:
	go run main.go add "userA" "working" "smile"

build:
	go build -o snulog main.go

lint:
	golangci-lint run

test:
	go test -v ./...


run-server:
	go run server/server.go

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO_SRC)

migrate-up:
	docker compose run --rm migrate

down:
	docker compose down

seed:
	docker exec -i snulog-db psql -U postgres -d snulog < db/seeds/dev_seed.sql

