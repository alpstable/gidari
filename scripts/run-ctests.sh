#!/bin/bash

docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d \
	--remove-orphans \
	--force-recreate \
	--build ctests

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run --rm ctests

