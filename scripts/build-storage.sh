#!/bin/sh

# remove volumes
rm -rf .db

# grant permissions
chmod +rwx ./third_party/docker/*.sh

# drop existing containers
docker compose -f "third_party/docker/storage.docker-compose.yaml" down

# prune containers
docker system prune --force

# re-create the storage containers
docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d --force-recreate --build mongo-coinbasepro
docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d --force-recreate --build mongo-coinbasepro2
docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d --force-recreate --build mongo-coinbasepro3
docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d --force-recreate --build postgres-coinbasepro
docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d --force-recreate --build postgres-polygon
docker-compose -f "third_party/docker/storage.docker-compose.yaml" up -d --force-recreate --build cache
