#!/bin/sh

set -e

sudo apt update
sudo apt remove -y golang-go

if [ $(command -v go &> /dev/null) ]; then
    sudo curl -OL https://golang.org/dl/go1.19.3.linux-amd64.tar.gz
    sudo tar -C /usr/local -xvf go1.19.3.linux-amd64.tar.gz
fi

export PATH=/usr/local/go/bin:"$PATH"

DIR=$(mktemp -d)

git clone https://github.com/logyball/quote-sender "$DIR"

cd "$DIR"

# build binary
sudo cp /home/"$(whoami)"/.env .
make build

sudo mv ./dist/quoteCats /home/"$(whoami)"/quoteCats

# add to cron
(crontab -l; echo "0 14 * * * /home/$(whoami)/quoteCats > /home/$(whoami)/cat.log 2>&1") | sort -u | crontab -

sudo rm -rf "$DIR"
