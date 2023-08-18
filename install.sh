#!/bin/sh
URL=$(curl -s https://api.github.com/repos/rafalb8/udns/releases/latest \
    | grep "browser_download_url.*amd64.tar.gz" \
    | cut -d : -f 2,3 \
    | tr -d \")

mkdir -p /tmp/udns
curl -sfL ${URL} | tar xvz -C /tmp/udns
sudo cp /tmp/udns/uDNS /bin/uDNS
rm -rf /tmp/udns