#!/bin/bash

# exit when any command fails
set -e

MDB_GW=$(docker inspect docker-mongo-1  -f "{{range .NetworkSettings.Networks }}{{.Gateway}}{{end}}")
PG_CBP_GW=$(docker inspect docker-postgres-coinbasepro-1  -f "{{range .NetworkSettings.Networks }}{{.Gateway}}{{end}}")
PG_P_GW=$(docker inspect docker-postgres-polygon-1  -f "{{range .NetworkSettings.Networks }}{{.Gateway}}{{end}}")

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-mongo-coinbasepro \
	--changelog-file=/changelog/changelog.xml \
	--headless=true \
	--url=mongodb://$MDB_GW:27017/coinbasepro \
	--log-level=debug update

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-mongo-polygon \
	--changelog-file=/changelog/changelog.xml \
	--headless=true \
	--url=mongodb://$MDB_GW:27017/polygon \
	--log-level=debug update

docker-compose -f "third_party/docker/storage.docker-compose.yaml" run liquibase-postgres-coinbasepro \
	--changelog-file=/changelog/changelog.xml \
	--url=jdbc:postgresql://$PG_CBP_GW:5432/coinbasepro \
	--log-level=debug \
	--username=postgres \
	--driver=org.postgresql.Driver update

# prune containers
docker system prune --force
