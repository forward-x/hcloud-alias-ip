package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"log"
	"net"
)

var client *hcloud.Client

func main() {
	apiToken := flag.String("token", "HCloud API Token", "-token <API Token>")
	strAliasIP := flag.String("alias-ip", "Alias IP", "-alias-ip <Alias-IP>")
	networkName := flag.String("network-name", "Network Name", "-network-name <Network Name>")
	serversLabel := flag.String("servers-label", "Servers Label", "-servers-label <Servers Label>")
	flag.Parse()

	if *apiToken == "" {
		log.Fatalf("No API Token specified!")
	}
	if *strAliasIP == "" {
		log.Fatalf("No Alias-IP specified!")
	}

	aliasIP := net.ParseIP(*strAliasIP)

	// Get System Hostname
	//serverName, err := os.Hostname()
	//if err != nil {
	//	panic(err)
	//}

	serverName := "test-db01"

	client = hcloud.NewClient(hcloud.WithToken(*apiToken))

	targetNetwork, _, err := client.Network.GetByName(context.Background(), *networkName)
	if err != nil {
		panic(err)
	}

	servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: *serversLabel,
		},
	})
	if err != nil {
		panic(err)
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
		log.Fatalf("Cannot remove Alias-IP '%s' from Server '%s', cannot find required network", aliasIP, server.Name)
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
		log.Fatalf("Cannot remove Alias-IP '%s' from Server '%s': %s", aliasIP, server.Name, err)
	}

	fmt.Printf("Alias-IP '%s' was removed from Server '%s'\n", aliasIP, server.Name)
}

func assignAliasIP(targetServerName string, targetNetwork *hcloud.Network, aliasIP net.IP) {
	targetServer, _, err := client.Server.GetByName(context.Background(), targetServerName)
	if err != nil {
		panic(err)
	}

	serverNet := findNetwork(targetServer.PrivateNet, *targetNetwork)

	serverAliases := append(serverNet.Aliases, aliasIP)

	_, _, err = client.Server.ChangeAliasIPs(context.Background(), targetServer, hcloud.ServerChangeAliasIPsOpts{
		Network:  serverNet.Network,
		AliasIPs: serverAliases,
	})

	if err != nil {
		log.Fatalf("Cannot assign Alias-IP '%s' to Server '%s'", aliasIP, targetServer.Name)
	}

	fmt.Printf("Alias-IP '%s' was assigned to Server '%s'\n", aliasIP, targetServer.Name)
}
