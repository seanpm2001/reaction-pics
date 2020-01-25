#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$DIR"/.. || exit 1

# Update repository
git checkout master
git fetch -tp
git pull

# Build and start container
docker pull "$(grep FROM Dockerfile | awk '{print $2}')"
docker build -t reaction-pics:production .
docker network inspect "reaction-pics" &>/dev/null ||
    docker network create --driver bridge "reaction-pics"
docker stop reaction-pics || true
docker container prune --force --filter "until=336h"
docker container rm reaction-pics || true
docker run \
    --detach \
    --restart=always \
    --publish="127.0.0.1:5003:5003" \
    --network="reaction-pics" \
    --name=reaction-pics reaction-pics:production

# Cleanup docker
docker image prune --force --filter "until=336h"

# Update nginx
sudo service nginx reload
