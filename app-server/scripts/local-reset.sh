#!/bin/bash

set -e

docker compose -f docker-compose.local.yaml down
docker volume prune -f
docker compose -f docker-compose.local.yaml up -d