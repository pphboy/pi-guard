package node

import (
	"context"
	"fmt"
	"go-node/models"
	"go-node/modules/nodegrpc"
	"go-node/modules/rproxy"
	ns "go-node/modules/service"
	gs "go-node/service"
	"go-node/sys"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"pglib"
	"pglib/cdns"
	pcent "pglib/center"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeBoot interface {
	Init()
	Install() error
	initService() error
	initGrpcServer() error
	initLog() error
	GetServiceManager() ns.ServiceManager
}

type NodeBootImpl struct {
	pkgService     gs.PkgService
	serviceManager ns.ServiceManager
	baseService    gs.BaseService
	reverseProxy   rproxy.ReverseProxy
	sysService     gs.SysService
	port           int
	nodeIp         string
	grpcPort       int
	nodeName       string
	center         string
	dnsM           cdns.CdnsManager
	bus            EventBus.Bus
}

func NewNode(center string, bs gs.BaseService, bus EventBus.Bus,
	port int, gp int, nodeName string,
	netseg string, dnsM cdns.CdnsManager) NodeBoot {

	n := &NodeBootImpl{
		bus:            bus,
		center:         center,
		baseService:    bs,
		pkgService:     gs.NewPkgService(bs, bus),
		serviceManager: ns.NewAppDirector(),
		reverseProxy:   rproxy.NewRProxyer(bs, bus),
		grpcPort:       gp,
		port:           port,
		nodeName:       nodeName,
		sysService:     gs.NewSysService(bs),
		dnsM:           dnsM,
	}

	bus.SubscribeAsync(ns.EVENT_ADDAPP, func(adj any) {
		nd, ok := adj.(*models.NodeApp)
		if ok {
			logrus.Infof("recived cmd, running, %+v", nd)
			appPath := path.Join(sys.PgSite(sys.PG_APP).Path, nd.NodeAppName)
			cfg, err := pglib.LoadPackageConfig(path.Join(appPath, pglib.PKGFILE_NAME))
			if err != nil {
				logrus.Errorf("load pkg.toml file error: cfg: %+v, err: %v", cfg, err)
				return
			}

			execFile := path.Join(appPath, cfg.Exec)
			if !pglib.IsFileExist(execFile) {
				logrus.Errorf("exec file not found in %s: cfg: %+v, err: %v", appPath, cfg, err)
				return
			}

			n.serviceManager.AddService(ns.NewRunnerApp(
				exec.Command(execFile),
				nd,
				n.baseService.DB,
				n.dnsM,
				n.nodeIp,
			))

			// 刷新反向代理缓存
			if err := n.reverseProxy.RefreshCache(); err != nil {
				logrus.Error("refresh reverseProxy cache,", err)
			}

		} else {
			logrus.Error("excp nodeinfo", nd)
		}
	}, false)

	// 获取结点IP，用于dns映射这个ip
	nt := gs.NewNodeNeter()
	nodeIp, err := nt.Ip4ByNetSegment(netseg)
	if err != nil {
		logrus.Fatal(err)
	} else if nt == nil {
		logrus.Fatal("get nil route,error network segment")
	}

	n.nodeIp = nodeIp.String()
	logrus.Info("node ipv4", n.nodeIp)

	// 安装node结点
	n.Install()

	return n
}

func ConnectCenter(center string) (pcent.CenterRadarClient, error) {
	logrus.Infof("try connect %s", center)
	conn, err := grpc.Dial(center, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pcent.NewCenterRadarClient(conn)
	return c, nil
}

func (n *NodeBootImpl) Init() {
	if err := n.initLog(); err != nil {
		logrus.Fatal("init log, ", err)
	}

	if err := n.initService(); err != nil {
		logrus.Fatal("init service", err)
	}

	go func() {
		logrus.Info("running grpc server")
		if err := n.initGrpcServer(); err != nil {
			logrus.Fatal("init grpc server", err)
		}

	}()

	var client pcent.CenterRadarClient
	d := fmt.Sprintf("%s.%s", n.nodeName, sys.ROOT_DOMAIN)

	// 尝试连接center
	connCenter := func() {
		if len(n.center) == 0 {
			return
		}

		for {
			c, err := ConnectCenter(n.center)
			if err != nil {
				logrus.Warnf("did not connect %s: %v", n.center, err)
				time.Sleep(1 * time.Second)
				continue
			}
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

			// FIXME: need get ip of this node
			//  center use ip:port connect node
			//  bind domain to ip , add to dns server
			resp, err := c.SendMe(ctx, &pcent.NodeReaction{
				Port:   int32(n.grpcPort),
				Ip:     n.nodeIp,
				Domain: d,
			})
			if err != nil {
				logrus.Error(err)
				time.Sleep(1 * time.Second)
				continue
			}
			logrus.Infof("send me ok, node %s ,resp: %v\n", d, resp)
			// 连上之后，开一个协程，进行健康检测
			client = c
			return
		}
	}

	healSaving := func() {
		// 健康检测
		for {
			if client != nil {
				ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
				resp, err := client.GetCenter(ctx, &pcent.NodeReaction{
					Port:   int32(n.grpcPort),
					Ip:     n.nodeIp,
					Domain: d,
				})
				if err != nil {
					logrus.Error(err)
					// 如果挂了，就开始connCenter
					connCenter()
					time.Sleep(3 * time.Second)
				}
				if resp != nil {
					logrus.Info("check center health ", resp.ProjectInfo.Domain)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}

	go func() {
		healSaving()
	}()

	go connCenter()

	n.startReverseHttp()
}

// service manager
func (n *NodeBootImpl) initService() error {
	logrus.Info("init service")

	if err := n.loadAppInfo(); err != nil {
		return err
	}

	return nil
}

func (n *NodeBootImpl) loadAppInfo() error {
	apps, err := n.pkgService.LoadAppList()

	if err != nil {
		return err
	}

	var loadErr string
	for _, v := range apps {
		// load Config
		appPath := path.Join(sys.PgSite(sys.PG_APP).Path, v.NodeAppName)
		cfg, err := pglib.LoadPackageConfig(path.Join(appPath, pglib.PKGFILE_NAME))
		if err != nil {
			logrus.Errorf("load pkg.toml file error: cfg: %+v, err: %v", cfg, err)
			loadErr = fmt.Sprint(loadErr, "\n", err)
			continue
		}

		execFile := path.Join(appPath, cfg.Exec)
		if !pglib.IsFileExist(execFile) {
			logrus.Errorf("exec file not found in %s: cfg: %+v, err: %v", appPath, cfg, err)
			loadErr = fmt.Sprintf("%s\n%s: exec not found", loadErr, execFile)
			continue
		}

		// 添加时就会启动
		// path plus exec file ,get absolute file path
		n.serviceManager.AddService(ns.NewRunnerApp(
			exec.Command(execFile),
			v,
			n.baseService.DB,
			n.dnsM,
			n.nodeIp,
		))
	}

	return nil
}

func (n *NodeBootImpl) startReverseHttp() {
	logrus.Info("start reverse http")
	m := http.ServeMux{}

	m.HandleFunc("/", n.reverseProxy.ReverseHandle)

	s := http.Server{
		Handler: &m,
		Addr:    fmt.Sprintf(":%d", n.port),
	}

	if err := s.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}
}

// 初始化grpc服务
func (n *NodeBootImpl) initGrpcServer() error {
	defer func() {
		if a := recover(); a != nil {
			logrus.Infof("grpc server error: %v", a)
		}
	}()

	ns := nodegrpc.NewNodeRpcService(n.serviceManager, n.baseService, n.bus)
	nodegrpc.RunRpcServer(n.grpcPort, ns)
	return nil
}

// 初始化grpc服务
func (n *NodeBootImpl) Install() error {
	// 直接安装，无所谓，反正只能安装一次
	return n.sysService.Install(n.nodeName)
}

func (n *NodeBootImpl) GetServiceManager() ns.ServiceManager {
	return n.serviceManager
}

func (n *NodeBootImpl) initLog() error {
	return nil
	nm := filepath.Join(sys.PgSite(sys.PG_LOGS).Path, fmt.Sprintf("%s_sys.log", n.nodeName))
	logFile, err := os.OpenFile(nm, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	logrus.SetOutput(logFile)
	// above warn , can be log in file

	return nil
}
