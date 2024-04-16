package centers

import (
	"context"
	"fmt"
	"go-ctrl/models"
	"go-node/modules/node"
	gs "go-node/service"
	"log"
	pcent "pglib/center"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCenter(t *testing.T) {
	port := 8000
	domain := "c1.pi.g"

	ininter := gs.NewIniter("../../../node/fs_root")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		pid := 1
		NewCenter(&models.PiProject{
			ProjectId:     &pid,
			ProjectStatus: &pid,
			Domain:        "c1.pi.g",
			Port:          &port,
		})

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
	wg.Wait()
}

func TestConnectGrpcCenter(t *testing.T) {

	port := 8000
	domain := "c1.pi.g"

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		NewCenter(&models.PiProject{
			Domain: domain,
			Port:   &port,
		})

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
