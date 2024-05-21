package cdns

import (
	"context"
	"fmt"
	"pi_dns/server"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CdnsManager interface {
	// 添加Host
	AddHosts(server.Host) error
	// 删除Host
	DelHosts(server.Host) error
	// 获取Host
	GetHosts(server.Host) error
}

func NewDnsManager(ip string, grpcport string) CdnsManager {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", ip, grpcport), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect %v", err)
	}

	c := server.NewHostManagerClient(conn)

	return &cdnsManager{
		dnsClient: c,
	}
}

type cdnsManager struct {
	dnsClient server.HostManagerClient
}

func (c *cdnsManager) AddHosts(h server.Host) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	resp, err := c.dnsClient.AddHosts(ctx, &h)
	if err != nil {
		return err
	}

	logrus.Infof("AddHost Message %v", resp.GetCode())
	return nil
}

func (c *cdnsManager) DelHosts(h server.Host) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	resp, err := c.dnsClient.DelHosts(ctx, &h)
	if err != nil {
		return err
	}
	logrus.Infof("DelHost Message %v", resp.GetCode())
	return nil
}

func (c *cdnsManager) GetHosts(h server.Host) error {
	return fmt.Errorf("no implemented method")
}
