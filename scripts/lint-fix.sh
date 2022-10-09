#!/bin/bash

# Check if the golangci-lint binary is installed locally. If it is, then use it, otherwise use the docker image.
if hash golangci-lint 2>/dev/null;
then
	golangci-lint run --fix;
else
	docker-compose -f "docker-compose.yml" up -d \
               	--remove-orphans \
                --force-recreate \
                --build lint

	docker-compose -f docker-compose.yml run
fi
