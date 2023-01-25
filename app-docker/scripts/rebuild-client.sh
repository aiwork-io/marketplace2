#!/bin/bash

set -e
ssh-add ~/.ssh/id_ed25519_ec2
git pull
docker compose stop cronclient
docker compose rm -f cronclient
docker compose build --no-cache cronclient
docker compose up -d cronclient
docker compose logs -f cronclient