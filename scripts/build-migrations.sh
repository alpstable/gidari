#!/bin/sh

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-mongo-coinbasepro \
	--changelog-file=/changelog/changelog.xml \
	--headless=true \
	--url=mongodb://mongo:27017/coinbasepro \
	--log-level=debug update

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-mongo-polygon \
	--changelog-file=/changelog/changelog.xml \
	--headless=true \
	--url=mongodb://mongo:27017/polygon \
	--log-level=debug update

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-postgres-coinbasepro \
	--changelog-file=/changelog/changelog.xml \
	--url=jdbc:postgresql://postgres-coinbasepro:5432/coinbasepro \
	--log-level=debug \
	--username=postgres \
	--driver=org.postgresql.Driver update

# prune containers
docker system prune --force
