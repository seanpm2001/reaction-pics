#!/bin/bash

# Setup server
sudo hostnamectl set-hostname reaction.pics

# Install go
sudo apt-get update
sudo apt-get install golang-go

# Clone repository
mkdir -p ~/gocode/src/github.com/albertyw/
mkdir -p ~/gocode/bin/
mkdir -p ~/gocode/pkg/
git clone git@github.com:albertyw/reaction-pics ~/gocode/src/github.com/albertyw/reaction-pics
ln -s gocode/src/github.com/albertyw/reaction-pics/ reaction-pics

# Setup env
ln -s .env.prod ~/reaction-pics/.env

# Install nginx
sudo add-apt-repository ppa:nginx/stable
sudo apt-get update
sudo apt-get install -y nginx

# Configure nginx
sudo rm -r /etc/nginx/sites-available
sudo ln -s ~/gocode/src/github.com/albertyw/reaction-pics/config/sites-available/app reaction.pics-app
sudo ln -s ~/gocode/src/github.com/albertyw/reaction-pics/config/sites-available/headers reaction.pics-headers
sudo rm -r /var/www/html

# Secure nginx
sudo mkdir /etc/nginx/ssl
sudo openssl dhparam -out /etc/nginx/ssl/dhparams.pem 2048
# Copy server.key and server.pem to /etc/nginx/ssl
sudo service nginx restart

# Compile
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
cd ~/gocode/src/github.com/albertyw/reaction-pics
dep ensure
make bins

# Set up go server service
sudo systemctl enable /home/ubuntu/gocode/src/github.com/albertyw/reaction-pics/config/reaction-pics.service
sudo systemctl start reaction-pics.service
