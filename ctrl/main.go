package main

import (
	"flag"
	"fmt"
	"go-ctrl/db"
	"go-ctrl/http"
	"go-ctrl/modules/admin"
	"go-ctrl/modules/appm"
	centers "go-ctrl/modules/center"
	"go-ctrl/modules/scripter"
	"go-ctrl/modules/ssher"
	"go-node/service"
	"os"
	"path"
	"pglib/cdns"
	"pi_dns/server"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	dataPath       = flag.String("path", "./", "store data path")
	ctrlName       = flag.String("name", "ctrl", "controller name")
	httpPort       = flag.Int("httpPort", 7431, "program running port")
	netseg         = flag.String("netseg", "192.168.56", "network segment")
	dnsManager     = flag.String("dnsManager", "192.168.56.104", "grpc dns server ip")
	dnsManagerPort = flag.String("dnsManagerPort", "50051", "grpc dns server port")
	picPath        = path.Join(*dataPath, "static")
	ctrlDomain     = flag.String("ctrlDomain", "ctrl.pi.g", "control domain")
)

func main() {
	flag.Parse()
	initOperation()

	logrus.Println("start running")
	s := http.NewCtrlHttp(*httpPort)

	s.RouterGroup("").Static("/assets", picPath)

	cd := cdns.NewDnsManager(*dnsManager, *dnsManagerPort)

	nt := service.NewNodeNeter()
	nodeIp, err := nt.Ip4ByNetSegment(*netseg)
	if err != nil {
		logrus.Fatal(err)
	} else if nt == nil {
		logrus.Fatal("get nil route,error network segment")
	}
	cd.AddHosts(server.Host{
		Domain: "ctrl.pi.g",
		Ips:    []string{nodeIp.String()},
	})

	// center 中心
	cm := centers.NewManager(s.RouterGroup("project"), cd, *netseg)

	// Project中心
	pm := centers.NewProjectManager(cm)

	// appStore
	appm.NewAppHttp(picPath, fmt.Sprintf("%s:%d", *ctrlDomain, *httpPort), s.RouterGroup("appm"))

	centers.NewProjectHttp(s.RouterGroup("project"), pm)

	// s.RouterGroup("project")
	ssher.NewSsherHttp(s.RouterGroup("term"))

	scripter.NewScripterHttp(s.RouterGroup("script"))

	admin.NewHttp(s.RouterGroup("user"))
	if err := s.Run(); err != nil {
		logrus.Fatal("ctrl http,", err)
	}
}

func initOperation() {
	// static file dir
	_, err := os.Stat(picPath)
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(picPath, 0744); err != nil {
			logrus.Error("mkdir picPath", picPath, err)
		}
	}

	logrus.Print("running main init")
	db.Init(*dataPath, *ctrlName)
}
