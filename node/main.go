package main

import (
	"flag"
	"go-node/modules/node"
	"go-node/service"
	"go-node/tool"
	"pglib/cdns"
	"strings"

	"github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
)

var (
	rootDir        = flag.String("root", "./", "root data dir")
	netseg         = flag.String("netseg", "192.168.56", "network segment")
	port           = flag.Int("port", 80, "http server port")
	dnsManager     = flag.String("dnsManager", "192.168.56.104", "grpc dns server ip")
	dnsManagerPort = flag.String("dnsManagerPort", "50051", "grpc dns server port")
	nodeName       = flag.String("name", tool.GetUUIDUpper(), "node name,default random uuid")
	grpcPort       = flag.Int("grpcPort", 9981, "grpc server port")
	center         = flag.String("center", "", "service center")
)

func main() {
	flag.Parse()

	logrus.Println("get param:\nrootDir:", *rootDir)
	// logrus.
	logrus.Println("http port:", *port)
	logrus.Println("grpc port:", *grpcPort)
	bus := EventBus.New()

	// 可能还需要读一下这个配置？
	// 先不管重启，能启动就行，重启的思路就是一个模块管理器，正在运动的模块关掉然后重新打开
	// 不过我觉得不太可取，还是要写一个类似于进程管理器的工具

	initer := service.NewIniter(*rootDir)
	dnsM := cdns.NewDnsManager(*dnsManager, *dnsManagerPort)
	// 如果安装时未填写相关的node名，则直接会以UUID的形式
	n := node.NewNode(*center,
		service.BaseService{DB: initer.GetDB()},
		bus,
		*port,
		*grpcPort,
		*nodeName,
		*netseg,
		dnsM)

	n.Init()
}

func checkNetSeg(netseg string) {
	ns := strings.Split(netseg, ".")
	if len(ns) != 3 {
		panic("error network segment")
	}
}
