#!/bin/sh

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-mongo-coinbasepro \
	--changelog-file=/changelog/changelog.xml \
	--headless=true \
	--url=mongodb://mongo-coinbasepro:27017/coinbasepro \
	--log-level=debug update

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-postgres-coinbasepro \
	--changelog-file=/changelog/changelog.xml \
	--url=jdbc:postgresql://postgres-coinbasepro:5432/coinbasepro \
	--log-level=debug \
	--username=postgres \
	--driver=org.postgresql.Driver update

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-postgres-polygon \
	--changelog-file=/changelog/changelog.xml \
	--url=jdbc:postgresql://postgres-polygon:5432/polygon \
	--log-level=debug \
	--username=postgres \
	--driver=org.postgresql.Driver update

# prune containers
docker system prune --force
