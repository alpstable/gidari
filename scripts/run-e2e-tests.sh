#!/bin/bash

set -e

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build e2e

docker-compose -f "docker-compose.yml" run --rm e2e
