#!/bin/bash

set -e

sudo chown -R $USER .db/mongo1
sudo chown -R $USER .db/mongo2
sudo chown -R $USER .db/mongo3
sudo chown -R $USER .db/postgres1

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build ctests

docker-compose -f "docker-compose.yml" run --rm \
	-e "CBP_KEY=$CBP_KEY" \
	-e "CBP_SECRET=$CBP_SECRET" \
	-e "CBP_PASSPHRASE=$CBP_PASSPHRASE" \
	-e "POL_BEARER_TOKEN=$POL_BEARER_TOKEN" \
	-e "MEEPMAN=$MEEPMAN" \
	ctests

