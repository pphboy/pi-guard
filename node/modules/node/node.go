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

type Node interface {
	Init()
	initService() error
	initGrpcServer() error
}

type NodeImpl struct {
	pkgService     gs.PkgService
	serviceManager ns.ServiceManager
	baseService    gs.BaseService
	reverseProxy   rproxy.ReverseProxy
	port           string
}

func NewNode(bs gs.BaseService, port string) Node {
	return &NodeImpl{
		baseService:    bs,
		pkgService:     gs.NewPkgService(bs),
		serviceManager: ns.NewAppDirector(),
		reverseProxy:   rproxy.NewRProxyer(bs),
		port:           port,
	}
}

func (n *NodeImpl) Init() {
	if err := n.initService(); err != nil {
		logrus.Fatal("init service", err)
	}
	if err := n.initGrpcServer(); err != nil {
		logrus.Fatal("init grpc server", err)
	}

	n.startReverseHttp()
}

// service manager
func (n *NodeImpl) initService() error {
	logrus.Info("init service")

	if err := n.loadAppInfo(); err != nil {
		return err
	}

	return nil
}

func (n *NodeImpl) loadAppInfo() error {
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

func (n *NodeImpl) startReverseHttp() {
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
func (n *NodeImpl) initGrpcServer() error {

	return nil
}
