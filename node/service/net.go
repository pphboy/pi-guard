package service

import (
	"net"
	"os"
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
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	for _, v := range ips {
		ip4 := v.To4()
		if ip4 != nil && strings.HasPrefix(ip4.String(), netseg) {
			return ip4, nil
		}
	}
	return nil, nil
}
