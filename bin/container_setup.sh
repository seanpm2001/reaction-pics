#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd $DIR/..

# Set locale
sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && locale-gen

# Install go and git
apt-get update
apt-get install -y golang-go git

# Setup env
rm -rf .env
ln -s .env.prod .env

# Compile
mkdir /root/gocode/bin
mkdir /root/gocode/pkg
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
export PATH="$PATH:$GOPATH/bin"
dep ensure
make bins
