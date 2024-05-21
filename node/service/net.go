package service

import (
	"net"
	"strings"
)

type NodeNet interface {
	Ip4ByNetSegment(netseg string) (net.IP, error)
}

func NewNodeNeter() NodeNet {
	return &nodeNeter{}
}

type nodeNeter struct {
}

func (n *nodeNeter) Ip4ByNetSegment(netseg string) (ipp net.IP, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, v := range addrs {
		ip4 := v.(*net.IPNet).IP
		if ip4 != nil && strings.HasPrefix(ip4.String(), netseg) {
			return ip4, nil
		}
	}
	return nil, nil
}
