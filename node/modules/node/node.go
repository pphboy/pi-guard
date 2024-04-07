package node

import (
	"fmt"
	"go-node/modules/rproxy"
	ns "go-node/modules/service"
	gs "go-node/service"
	"go-node/sys"
	"net/http"
	"os/exec"
	"path"
	"pglib"

	"github.com/sirupsen/logrus"
)

type NodeBoot interface {
	Init()
	Install() error
	initService() error
	initGrpcServer() error
}

type NodeBootImpl struct {
	pkgService     gs.PkgService
	serviceManager ns.ServiceManager
	baseService    gs.BaseService
	reverseProxy   rproxy.ReverseProxy
	sysService     gs.SysService
	port           string
	grpcPort       string
	nodeName       string
}

func NewNode(bs gs.BaseService, port string, gp string, nodeName string) NodeBoot {

	n := &NodeBootImpl{
		baseService:    bs,
		pkgService:     gs.NewPkgService(bs),
		serviceManager: ns.NewAppDirector(),
		reverseProxy:   rproxy.NewRProxyer(bs),
		grpcPort:       gp,
		port:           port,
		nodeName:       nodeName,
		sysService:     gs.NewSysService(bs),
	}

	// 安装node结点
	n.Install()

	return n
}

func (n *NodeBootImpl) Init() {
	if err := n.initService(); err != nil {
		logrus.Fatal("init service", err)
	}
	if err := n.initGrpcServer(); err != nil {
		logrus.Fatal("init grpc server", err)
	}

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
		Addr:    fmt.Sprintf(":%s", n.port),
	}

	s.ListenAndServe()
}

// 初始化grpc服务
func (n *NodeBootImpl) initGrpcServer() error {

	return nil
}

// 初始化grpc服务
func (n *NodeBootImpl) Install() error {
	// 直接安装，无所谓，反正只能安装一次
	return n.sysService.Install(n.nodeName)
}
