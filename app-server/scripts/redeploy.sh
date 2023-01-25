#!/bin/bash

set -e


cd ~/aiworkmarketplace
git pull

echo "AWS_ACCESS_KEY_ID: $AWS_ACCESS_KEY_ID"

docker compose build api wallet-helpers --no-cache
docker compose rm -f api wallet-helpers
docker compose up -d api wallet-helpers