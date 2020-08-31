#!/bin/bash

# This script will build and deploy a new docker image

set -exuo pipefail
IFS=$'\n\t'

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$DIR"/.. || exit 1

DEPLOY_BRANCH="${1:-}"
BRANCH="$(git rev-parse --abbrev-ref HEAD)"
set +x  # Do not print contents of .env
source .env
set -x

if [ -n "$DEPLOY_BRANCH" ]; then
    # Update repository
    git checkout "$DEPLOY_BRANCH"
    git fetch -tp
    git pull
fi

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

if [ "$ENV" = "production" ] && [ "$BRANCH" = "master" ]; then
    # Cleanup docker
    docker system prune --force --filter "until=168h"
    docker volume prune --force

    # Update nginx
    sudo service nginx reload
fi
