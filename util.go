package main

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"net"
)

func findNetwork(serverNets []hcloud.ServerPrivateNet, targetNet hcloud.Network) *hcloud.ServerPrivateNet {
	for _, serverNet := range serverNets {
		if serverNet.Network.ID == targetNet.ID {
			return &serverNet
		}
	}

	return nil
}

func indexOf(elements []net.IP, item net.IP) int {
	for index, element := range elements {
		if element.Equal(item) {
			return index
		}
	}

	return -1
}

func removeByIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
