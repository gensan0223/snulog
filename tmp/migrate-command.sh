#!/bin/sh
migrate \
  -source=file:///migrations \
  -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@db:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" \
  -verbose up

