package centers

import (
	"context"
	"fmt"
	"go-ctrl/models"
	"net"
	"pglib/center"
	"snproto"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO:
// 实现Node客户端
//    - 通过Node客户端识别结点是否是Center结点
//    - Center结点上就可以拿到所有Node的信息
//    - 这个Center也可以独占一个端口，总之就是一个程序
//    - 目前直接放到Ctrl里面吧，建立一个项目就会开启一个Center线程
//    - center线程然后会去扫描其自身的Node结点
//    - Center就是一个单独的服务

/*
Center功能：
    - NodeList
		- 获取Node之后，然后就通过grpcClient连到Node的GrpcServer
		- 通过Node的rpc服务拿到Node相关的信息
	- 暴露Center中心
	- Center Grpc Server
*/

// TODO:
// Center写完后，马上写云端App，这样就有办法将云端APP装到Center上面。
// 顺序：
// 		- Center
// 		- NodeList
// 		- 监控，监控应该要放到Node单独的界面，而不是独立一个界面
// 		- 云端App
// 		- ssher
// 		- scriper

type Center interface {
	RegisterNode(*center.NodeReaction) error
	RunRpcServer() error
	Info() *models.PiProject
	Port() int32
	Domain() string
}

type centerServer struct {
	center.UnimplementedCenterRadarServer
	center Center
}

func (c *centerServer) SendMe(ctx context.Context, nr *center.NodeReaction) (
	*center.CenterReaction, error) {
	p := c.center.Info()
	logrus.Infof("register %s:%d node reaction", nr.Domain, nr.Port)
	c.center.RegisterNode(nr)
	return &center.CenterReaction{
		ProjectInfo: p.Msg(),
	}, nil
}

type centerImpl struct {
	project  *models.PiProject
	nodeList []*NodeGrpcClient
}

type NodeGrpcClient struct {
	snproto.MonitorServiceClient
	snproto.NodeAppServiceClient
	snproto.NodeServiceClient
}

func NewCenter(p *models.PiProject) Center {
	ci := &centerImpl{
		project: p,
	}

	go ci.RunRpcServer()
	return ci
}

func (c *centerImpl) RegisterNode(nr *center.NodeReaction) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", nr.Domain, nr.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	ngc := &NodeGrpcClient{
		snproto.NewMonitorServiceClient(conn),
		snproto.NewNodeAppServiceClient(conn),
		snproto.NewNodeServiceClient(conn),
	}

	// 连接Node的Grpc，然后将连接存起来
	// 定时拿一下性能信息，类似于一个心跳
	c.nodeList = append(c.nodeList, ngc)
	// timeout 1 minute
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	ndInfo, err := ngc.GetNodeInfo(ctx, &snproto.Empty{})
	if err != nil {
		return err
	}

	logrus.Infof("center %s connect [node:%s] succeed!", c.project.ProjectName, ndInfo.NodeDomain)
	return nil
}

func (c *centerImpl) RunRpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port()))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	center.RegisterCenterRadarServer(s, &centerServer{
		center: c,
	})

	logrus.Printf("ceneter %s server listening at %v", c.project.Domain, lis.Addr())

	if err := s.Serve(lis); err != nil {
		logrus.Errorf("%v %s, grpc shutdown",
			c.project.ProjectId, c.project.ProjectName)
		return err
	}

	return nil
}

func (c *centerImpl) Info() *models.PiProject {
	return c.project
}

func (c *centerImpl) Port() int32 {
	return int32(*c.project.Port)

}

func (c *centerImpl) Domain() string {
	return c.project.Domain
}
