#!/bin/bash

docker-compose -f "docker-compose.yml" up -d \
    --remove-orphans \
    --force-recreate \
    --build fmt

docker-compose -f docker-compose.yml run --rm fmt
