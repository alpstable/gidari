#!/bin/sh

# remove volumes
rm -rf .db

# grant permissions
chmod +rwx ./third_party/docker/*.sh

# drop existing containers
docker compose -f "third_party/docker/storage.docker-compose.yaml" down

# prune containers
docker system prune --force

docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d \
	--remove-orphans \
	--force-recreate \
	--build mongo

docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d \
	--remove-orphans \
	--force-recreate \
	--build mongo2

docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d \
	--remove-orphans \
	--force-recreate \
	--build mongo3

docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d \
	--remove-orphans \
	--force-recreate \
	--build postgres-coinbasepro

docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d \
	--remove-orphans \
	--force-recreate \
	--build postgres-polygon
