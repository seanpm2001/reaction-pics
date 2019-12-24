#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$DIR/.." || exit 1

# Set locale
sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && locale-gen

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

# Compile code
npm install
npm run minify
npm prune --production
