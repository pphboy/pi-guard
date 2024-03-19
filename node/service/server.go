package service

import (
	"fmt"
	"go-node/models"
	"go-node/sys"
	"path"
	"pglib"

	"github.com/sirupsen/logrus"
)

type Server interface {
	init()
	InitRunApps([]*models.NodeApp) error
}

type nodeServerImpl struct {
	pkgService PkgService
}

func (n *nodeServerImpl) init() {

}

// 启动服务
func (n *nodeServerImpl) InitRunApps() error {
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
		if !IsFileExist(execFile) {
			logrus.Errorf("exec file not found in %s: cfg: %+v, err: %v", appPath, cfg, err)
			loadErr = fmt.Sprintf("%s\n%s: exec not found", loadErr, execFile)
			continue
		}

		// path plus exec file ,get absolute file path
	}

	return nil
}
