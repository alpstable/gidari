#!/bin/bash

set -e

docker-compose -f "docker-compose.yml" up -d --build integration
docker-compose -f "docker-compose.yml" run --rm integration -tags=$1 -count=$2
