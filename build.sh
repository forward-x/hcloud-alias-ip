#!/bin/bash

GOARCH=arm64 GOOS=linux go build -o app *.go

# Tar the executable, service, and install script
mkdir -p hcloud-alias-ip
cp app hcloud-alias-ip/

tar -czvf hcloud-alias-ip.arm64.linux.tar.gz hcloud-alias-ip
rm -rf hcloud-alias-ip