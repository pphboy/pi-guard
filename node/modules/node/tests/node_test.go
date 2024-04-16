package tests

import (
	"go-node/modules/node"
	gs "go-node/service"
	"testing"
	"time"
)

func TestNodeStart(t *testing.T) {
	ininter := gs.NewIniter("../../../fs_root")

	n := node.NewNode("", gs.BaseService{DB: ininter.GetDB()}, 8090, 9981, "ndev")

	n.Init()
	time.Sleep(29 * time.Second)
}
