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
docker run --detach --restart always -p 127.0.0.1:5003:5003 \
    --mount type=bind,source="$(pwd)"/tumblr/data,target=/root/gocode/src/github.com/albertyw/reaction-pics/tumblr/data \
    --name reaction-pics reaction-pics:production

# Cleanup docker
docker image prune --force --filter "until=336h"

# Update nginx
sudo service nginx reload
