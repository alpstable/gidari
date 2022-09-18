#!/bin/bash

set -e

# remove volumes
rm -rf .db

# remake .db
mkdir .db
chmod +rwx .db

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
docker-compose -f "docker-compose.yml" exec -T mongo1 bash -c '
mongosh <<EOF
var config = {
    "_id": "dbrs",
    "version": 1,
    "members": [
        {
            "_id": 1,
            "host": "mongo1:27017",
            "priority": 3
        },
        {
            "_id": 2,
            "host": "mongo2:27017",
            "priority": 2
        },
        {
            "_id": 3,
            "host": "mongo3:27017",
            "priority": 1
        }
    ]
};
rs.initiate(config, { force: true });
rs.status();
EOF
'

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build postgres1
