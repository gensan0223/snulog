PROTO_SRC=proto/logs.proto

run:
	go run main.go $(ARGS)

run-fetch:
	make run ARGS="fetch"

run-add:
	make run ARGS="add"

build:
	go build -o snulog main.go

run-server:
	go run server/server.go

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO_SRC)

up:
	docker compose up -d

migrate-up:
	docker compose run --rm migrate

down:
	docker compose down

seed:
	docker exec -i snulog-db psql -U postgres -d snulog < db/seeds/dev_seed.sql

