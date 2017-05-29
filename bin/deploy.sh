#!/bin/bash

# Update repository
cd ~/gocode/src/github.com/albertyw/reaction-pics/ || exit 1
git checkout master
git pull
glide install
go build

# Restart services
sudo service nginx restart
sudo systemctl restart reaction-pics.service
