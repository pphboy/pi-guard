package service

import (
	"go-node/service"
	"net"
	"testing"
)

func TestNodeNter(t *testing.T) {
	n := service.NewNodeNeter()
	ip, err := n.Ip4ByNetSegment("127.0.0")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip.String())
}

func TestIntefaces(t *testing.T) {
	i, err := net.InterfaceAddrs()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range i {
		t.Log(v.(*net.IPNet).IP.String())
	}
}
