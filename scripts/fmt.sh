#!/bin/bash

# Check if the "gofumpt" binary is installed locally. If it is, then use it, otherwise use the docker image.
if hash gofumpt 2>/dev/null;
then
	gofumpt -l -w .
	golangci-lint run --fix
else
	docker-compose -f "docker-compose.yml" up -d \
    	--remove-orphans \
    	--force-recreate \
    	--build fmt

	docker-compose -f docker-compose.yml run --rm fmt
fi
