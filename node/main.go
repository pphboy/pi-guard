package main

import (
	"flag"
	"go-node/modules/node"
	"go-node/service"
	"go-node/tool"

	"github.com/sirupsen/logrus"
)

var (
	rootDir  = flag.String("root", "./", "root data dir")
	port     = flag.Int("port", 80, "http server port")
	nodeName = flag.String("name", tool.GetUUIDUpper(), "node name,default random uuid")
	grpcPort = flag.Int("grpcPort", 9981, "grpc server port")
)

func main() {
	logrus.Println("get param:\nrootDir:", *rootDir)
	// logrus.
	logrus.Println("http port:", *port)
	logrus.Println("grpc port:", *grpcPort)

	// 可能还需要读一下这个配置？
	// 先不管重启，能启动就行，重启的思路就是一个模块管理器，正在运动的模块关掉然后重新打开
	// 不过我觉得不太可取，还是要写一个类似于进程管理器的工具

	initer := service.NewIniter(*rootDir)
	// 如果安装时未填写相关的node名，则直接会以UUID的形式
	n := node.NewNode(service.BaseService{DB: initer.GetDB()}, *port, *grpcPort, *nodeName)

	n.Init()
}
