#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

# Setup server
sudo hostnamectl set-hostname reaction.pics

# Clone repository
cd ~ || exit 1
git clone git@github.com:albertyw/reaction-pics

# Install nginx
sudo add-apt-repository ppa:nginx/stable
sudo apt-get update
sudo apt-get install -y nginx

# Configure nginx
sudo rm -r /etc/nginx/sites-available
sudo rm -r /var/www/html

# Secure nginx
sudo mkdir /etc/nginx/ssl
curl https://ssl-config.mozilla.org/ffdhe2048.txt | sudo tee /etc/nginx/ssl/dhparams.pem > /dev/null
# Copy server.key and server.pem to /etc/nginx/ssl
sudo service nginx restart

# Set up docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce
sudo usermod -aG docker "${USER}"
