#!/bin/bash

set -e

# remove volumes
rm -rf .db

# drop existing containers
docker compose -f "docker-compose.yml" down

# prune containers
docker system prune --force

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build mongo1

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build mongo2

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build mongo3

echo "Waiting for MongoDB to start..."
sleep 5

echo "Creating replica set..."
docker-compose -f "docker-compose.yml" exec -T mongo1 /scripts/rs-init.sh

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build postgres1
