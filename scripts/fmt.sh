#!/bin/bash

if hash gofumpt 2>/dev/null;
then
	gofumpt -l -w .
else

docker-compose -f "docker-compose.yml" up -d \
    --remove-orphans \
    --force-recreate \
    --build fmt

docker-compose -f docker-compose.yml run --rm fmt

fi
