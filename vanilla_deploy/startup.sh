#!/bin/sh

set -e

if [ "$(echo $PATH | grep -c go)" = 0 ]; then
    export PATH=/usr/local/go/bin:/usr/local/go:"$PATH"
fi

# install go if necessary
if [ $(command -v go > /dev/null 2>&1) ]; then
    sudo apt update
    sudo apt remove -y golang-go

    if [ $(command -v go &> /dev/null) ]; then
        sudo curl -OL https://golang.org/dl/go1.19.3.linux-amd64.tar.gz
        sudo tar -C /usr/local -xvf go1.19.3.linux-amd64.tar.gz
    fi

    export PATH=/usr/local/go/bin:"$PATH"
fi

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

# install promtail if necessary
if [ $(command -v /opt/promtail/promtail > /dev/null 2>&1) ]; then
    echo "promtail not found, installing and starting as service"

    sudo rm -rf /opt/promtail

    PROMTAIL_VERSION=$(curl -s "https://api.github.com/repos/grafana/loki/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')
    sudo mkdir -p /opt/promtail
    sudo wget -qO /opt/promtail/promtail.gz "https://github.com/grafana/loki/releases/download/v${PROMTAIL_VERSION}/promtail-linux-amd64.zip"
    sudo gunzip /opt/promtail/promtail.gz
    sudo chmod a+x /opt/promtail/promtail

    sudo mv /home/"$(whoami)"/promtail.service /etc/systemd/system/promtail.service
    sudo mv /home/"$(whoami)"/promtail-config.yaml /opt/promtail/promtail-config.yaml


    sudo systemctl daemon-reload
    sudo systemctl start promtail.service
fi
