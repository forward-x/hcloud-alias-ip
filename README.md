# Hetzner Cloud Alias IP Assigner
A CLI utility to assign an Alias-IP to the Server executing it.

## Requirements
+ Hostname equals Cloud Server name
+ Read/Write Cloud API Key
+ Assigned Label to Servers for organizing pool of servers the Alias-IP should belong
+ Network Name the Alias-IP should belong

## Usage:

```bash
hcloud-alias-ip -alias-ip <Alias-IP> \
                -network-name <Network Name> \
                -server-label <Server Label> \
                -token <API Token>
```