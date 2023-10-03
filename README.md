# Hetzner Cloud Alias IP Assigner
A CLI utility to assign an Alias-IP to the Server executing it.

## Requirements
+ Hostname equals Cloud Server name
+ Read/Write Cloud API Key
+ Assigned Label to Servers for organizing pool of servers the Alias-IP should belong
+ Network Name the Alias-IP should belong

## Installing

Download and unpack `hcloud-alias-ip`

```bash
HCLOUD_ALIAS_IP_VERSION=v0.0.3
wget -O hcloud-alias-ip.tar.gz https://github.com/glav-kod/hcloud-alias-ip/releases/download/$HCLOUD_ALIAS_IP_VERSION/hcloud-alias-ip-$HCLOUD_ALIAS_IP_VERSION-linux-amd64.tar.gz
tar -zxvf hcloud-alias-ip.tar.gz
```

Move it to `/usr/local/bin` directory
```bash
mv hcloud-alias-ip /usr/local/bin
```

## Usage

```bash
hcloud-alias-ip -alias-ip <Alias-IP> \
                -network-name <Network Name> \
                -server-label <Server Label> \
                -token <API Token>
```