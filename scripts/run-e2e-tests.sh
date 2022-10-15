#!/bin/bash

set -e

docker-compose -f "docker-compose.yml" up -d --build e2e
docker-compose -f "docker-compose.yml" run --rm e2e
