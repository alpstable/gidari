#!/bin/bash

set -e

# Try to export the environment if it exists.
[[ -f /etc/alpine-hodler/auth.env ]] && export $(cat /etc/alpine-hodler/auth.env | xargs)

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build ci

docker-compose -f "docker-compose.yml" run --rm \
	-e "CBP_KEY=$CBP_KEY" \
	-e "CBP_SECRET=$CBP_SECRET" \
	-e "CBP_PASSPHRASE=$CBP_PASSPHRASE" \
	-e "POL_BEARER_TOKEN=$POL_BEARER_TOKEN" \
	ci


