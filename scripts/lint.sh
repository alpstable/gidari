#!/bin/bash

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
        --force-recreate \
        --build lint

docker-compose -f docker-compose.yml run --rm lint
