#!/bin/bash

set -e

docker-compose -f "docker-compose.yml" up -d --build cmd
docker-compose -f "docker-compose.yml" run --rm cmd
