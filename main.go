package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"net"
	"os"
)

var client *hcloud.Client

func main() {
	apiToken := flag.String("token", "", "-token <API Token>")
	strAliasIP := flag.String("alias-ip", "", "-alias-ip <Alias-IP>")
	networkName := flag.String("network-name", "", "-network-name <Network Name>")
	serverLabel := flag.String("server-label", "", "-server-label <Server Label>")
	flag.Parse()

	if *apiToken == "" {
		ShowUsage()
	}
	if *strAliasIP == "" {
		ShowUsage()
	}
	if *networkName == "" {
		ShowUsage()
	}
	if *serverLabel == "" {
		ShowUsage()
	}

	aliasIP := net.ParseIP(*strAliasIP)

	serverName, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Cannot get hostname: %s", err))
	}

	client = hcloud.NewClient(hcloud.WithToken(*apiToken))

	targetNetwork, _, err := client.Network.GetByName(context.Background(), *networkName)
	if err != nil {
		panic(fmt.Sprintf("Cannot find network '%s': %s", *networkName, err))
	}

	servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: *serverLabel,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Cannot list servers by label '%s': %s", *serverLabel, err))
	}

	currentServer := findServerByAliasIP(servers, targetNetwork, aliasIP)

	if currentServer != nil {
		if currentServer.Name == serverName {
			fmt.Printf("Alias-IP '%s' is already assigned to Server '%s'", aliasIP, serverName)
			return
		}

		removeAliasIP(currentServer, targetNetwork, aliasIP)
	}

	assignAliasIP(serverName, targetNetwork, aliasIP)
}

func findServerByAliasIP(servers []*hcloud.Server, targetNetwork *hcloud.Network, aliasIP net.IP) *hcloud.Server {
	for _, server := range servers {
		serverNet := findNetwork(server.PrivateNet, *targetNetwork)
		if serverNet == nil {
			continue
		}

		index := indexOf(serverNet.Aliases, aliasIP)
		if index == -1 {
			continue
		}

		return server
	}

	return nil
}

func removeAliasIP(server *hcloud.Server, targetNetwork *hcloud.Network, aliasIP net.IP) {
	serverNet := findNetwork(server.PrivateNet, *targetNetwork)
	if serverNet == nil {
		panic(fmt.Sprintf("Cannot remove Alias-IP '%s' from Server '%s', cannot find required network", aliasIP, server.Name))
	}

	index := indexOf(serverNet.Aliases, aliasIP)
	if index == -1 {
		return
	}

	aliases := removeByIndex(serverNet.Aliases, index)

	_, _, err := client.Server.ChangeAliasIPs(context.Background(), server, hcloud.ServerChangeAliasIPsOpts{
		Network:  serverNet.Network,
		AliasIPs: aliases,
	})
	if err != nil {
		panic(fmt.Sprintf("Cannot remove Alias-IP '%s' from Server '%s': %s", aliasIP, server.Name, err))
	}

	fmt.Printf("Alias-IP '%s' was removed from Server '%s'\n", aliasIP, server.Name)
}

func assignAliasIP(targetServerName string, targetNetwork *hcloud.Network, aliasIP net.IP) {
	targetServer, _, err := client.Server.GetByName(context.Background(), targetServerName)
	if err != nil {
		panic(fmt.Sprintf("Cannot get Server '%s': %s", targetServerName, err))
	}

	serverNet := findNetwork(targetServer.PrivateNet, *targetNetwork)

	serverAliases := append(serverNet.Aliases, aliasIP)

	_, _, err = client.Server.ChangeAliasIPs(context.Background(), targetServer, hcloud.ServerChangeAliasIPsOpts{
		Network:  serverNet.Network,
		AliasIPs: serverAliases,
	})

	if err != nil {
		panic(fmt.Sprintf("Cannot assign Alias-IP '%s' to Server '%s'", aliasIP, targetServer.Name))
	}

	fmt.Printf("Alias-IP '%s' was assigned to Server '%s'\n", aliasIP, targetServer.Name)
}

func ShowUsage() {
	fmt.Printf(`Usage:
hcloud-alias-ip -alias-ip <Alias-IP> \
                -network-name <Network Name> \
                -server-label <Server Label> \
                -token <API Token>
`)
	os.Exit(-1)
}
