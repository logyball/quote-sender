#!/bin/sh

set -e

sudo apt update
sudo apt install golang-go

DIR=$(mktemp -d)

git clone https://github.com/logyball/quote-sender "$DIR"

cd "$DIR"

# build binary
make build
sudo mv ./dist/quoteCats /usr/local/bin/quoteCats

sudo rm -rf "$DIR"