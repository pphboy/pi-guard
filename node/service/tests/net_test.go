package service

import (
	"go-node/service"
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
