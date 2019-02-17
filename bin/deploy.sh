#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$DIR"/..

# Update repository
git checkout master
git fetch -tp
git pull

# Build and start container
docker build -t reaction-pics:production .
docker stop reaction-pics || echo
docker container prune --force --filter "until=336h"
docker container rm reaction-pics
docker run --detach --restart always -p 127.0.0.1:5003:5003 \
    --name reaction-pics reaction-pics:production

# Cleanup docker
docker image prune --force --filter "until=336h"

# Update nginx
sudo service nginx reload
