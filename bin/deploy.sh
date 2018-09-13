#!/bin/bash

# Update repository
cd ~/gocode/src/github.com/albertyw/reaction-pics/ || exit 1
git checkout master
git fetch -tp
git pull

# Build and start container
docker build -t reaction-pics:production .
docker stop reaction-pics || echo
docker container prune -f
docker run --detach --restart always -p 127.0.0.1:5003:5003 --name reaction-pics reaction-pics:production

# Cleanup docker
docker image prune -f

# Update nginx
sudo service nginx reload
