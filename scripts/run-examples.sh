#!/bin/bash

set -e

docker-compose -f "docker-compose.yml" up -d --build examples

# Iterate over all the directories in the examples/ directory.
for dir in examples/*/
do
	docker-compose -f "docker-compose.yml" run --rm examples /bin/bash -c "cd $dir && go run main.go"
done

