#!/bin/sh

set -e

sudo apt update
sudo apt remove -y golang-go

if [ $(command -v go &> /dev/null) ] then
    sudo curl -OL https://golang.org/dl/go1.19.3.linux-amd64.tar.gz
    sudo tar -C /usr/local -xvf go1.19.3.linux-amd64.tar.gz
fi

if [[ "$PATH" != *"/usr/local/go"* ]]; then
    export PATH=/usr/local/go/bin:"$PATH"
fi

DIR=$(mktemp -d)

git clone https://github.com/logyball/quote-sender "$DIR"

cd "$DIR"

# build binary
make build
sudo mv ./dist/quoteCats /usr/local/bin/quoteCats
sudo mv /home/"$(whoami)"/.env /usr/local/bin/.env

sudo rm -rf "$DIR"