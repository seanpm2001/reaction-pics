#!/bin/bash

# Update repository
cd ~/gocode/src/github.com/albertyw/reaction-pics/ || exit 1
git checkout master
git pull
dep ensure
go build

# Update permissions
sudo chmod -R 777 tumblr/data

# Restart services
sudo service nginx restart
sudo systemctl restart reaction-pics.service
