#!/usr/bin/env bash

docker-compose -f "meta.docker-compose.yaml" up -d --build --remove-orphans
# docker-compose -f "meta.docker-compose.yaml" run test_generate
docker-compose -f "meta.docker-compose.yaml" run generate
# go get -d github.com/99designs/gqlgen
# go run github.com/99designs/gqlgen generate
