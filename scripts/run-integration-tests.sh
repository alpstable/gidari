#!/bin/bash

set -e

NAME=$1

docker-compose -f "docker-compose.yml" up -d --build $NAME
docker-compose -f "docker-compose.yml" run --rm $NAME
