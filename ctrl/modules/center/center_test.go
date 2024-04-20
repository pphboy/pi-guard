package centers

import (
	"context"
	"fmt"
	"go-ctrl/db"
	"go-ctrl/models"
	"go-node/modules/node"
	gs "go-node/service"
	"log"
	pcent "pglib/center"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	db.Init("./", "test")
}

func TestCenter(t *testing.T) {

	g := gin.Default()
	go func() {
		g.Run(":9901")
	}()

	port := 8000
	domain := "c1.pi.g"

	ininter := gs.NewIniter("../../../node/fs_root")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		pid := 1
		NewCenter(nil, &models.PiProject{
			ProjectId:     &pid,
			ProjectStatus: &pid,
			Domain:        "c1.pi.g",
			Port:          &port,
		}, g.Group("/"))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 2)
		n := node.NewNode(fmt.Sprintf("%s:%d", domain, port), gs.BaseService{DB: ininter.GetDB()}, 8081, 9981, "ndev")
		n.Init()
		wg.Done()
	}()

	select {}
}

func TestConnectGrpcCenter(t *testing.T) {
	g := gin.Default()
	go func() {
		g.Run(":9901")
	}()

	port := 8000
	domain := "c1.pi.g"

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		NewCenter(nil, &models.PiProject{
			Domain: domain,
			Port:   &port,
		}, g.Group("project"))

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Second)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", domain, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c := pcent.NewCenterRadarClient(conn)
		c.SendMe(context.Background(), &pcent.NodeReaction{
			Port:   6666,
			Domain: "ndev.pi.g",
		})
		wg.Done()
	}()

	wg.Wait()
}
