package service

import (
	"errors"
	"fmt"
	"go-node/dao"
	"go-node/models"
	"go-node/sys"
	"go-node/tool"
	"io"
	"net/http"
	"os"
	"path/filepath"
	pglib "pglib"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrLoadAppNotExist = errors.New("load app not existing")
)

type PkgService interface {
	LoadAppList() (narr []*models.NodeApp, err error)
	UninstallApp(np *models.NodeApp) error
	InstallApp(pc *models.PiCloudApp) error
}

func NewPkgService(bs BaseService) PkgService {
	return &pkgServiceImpl{
		BaseService: bs,
		appDao: &dao.AppDaoImpl{
			Db: bs.DB,
		},
		sysDao: &dao.SysDaoImpl{
			Db: bs.DB,
		},
	}
}

type pkgServiceImpl struct {
	BaseService
	appDao dao.AppDao
	sysDao dao.SysDao
}

// 加载应用列表
func (p *pkgServiceImpl) LoadAppList() (narr []*models.NodeApp, err error) {
	// 加载 app数据库的记录即可
	// 然后去/pg/app文件中拿数据，看数据是否存在
	// 加载 /pg/app 路径中所有的应用
	// 随带去验证一下，目录下的app是否存在

	as, err := p.appDao.GetAll()
	if err != nil {
		return nil, err
	}
	var noExist string
	for _, v := range as {
		appDir := fmt.Sprintf("%s/%s", sys.PgSite(sys.PG_APP).Path, v.NodeAppName)
		if !IsDirExist(appDir) {
			noExist = fmt.Sprintf("%s\n[%v/%v, don't have exec file]",
				noExist, v.NodeAppId, v.NodeAppName)
			continue
		}
		narr = append(narr, v)
	}

	if len(noExist) > 0 {
		return narr, fmt.Errorf("%w, %s", ErrLoadAppNotExist, noExist)
	}

	return narr, nil
}

// 卸载应用
func (p *pkgServiceImpl) UninstallApp(np *models.NodeApp) error {
	// 关闭当前端口的应用
	// 删除 /pg/app 路径中的应用，移动到 pg/.trash
	oldPath := fmt.Sprintf("%s/%s", sys.PgSite(sys.PG_APP).Path, np.NodeAppName)
	newPath := fmt.Sprintf("%s/%s_%s", sys.PgSite(sys.PG_TRASH).Path, np.NodeAppName, tool.GetUUIDUpper())
	logrus.Printf("move to pg/.trash\noldPath:%s\nnewPath:%s\n", oldPath, newPath)
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("uninstall app, %w", err)
	}

	// 删除数据库中的记录，软删除
	if err := p.appDao.Delete(np); err != nil {
		// rename back
		os.Rename(newPath, oldPath)
		return err
	}

	return nil
}

// 安装应用
//   - 安装之前，需要判断一下版本，将原版本的删掉，改成新版本的
func (p *pkgServiceImpl) InstallApp(pc *models.PiCloudApp) error {
	existApp, err := p.appDao.GetByName(pc.AppName)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return fmt.Errorf("app dao get by name,%w", err)
	} else {
		if existApp.NodeAppVersion > pc.AppVersion {
			return fmt.Errorf("install app has lower version")
		}
	}

	resp, err := http.Get(pc.AppSite)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logrus.Println("pkg file size:", resp.ContentLength)

	// 下载到压缩包到 /tmp 临时路径,并解压，然后读取 config
	tempPkgFile, err := os.CreateTemp("/tmp", "pgdown")
	if err != nil {
		logrus.Error("failed create download temp file")
		return err
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read all file, %w", err)
	}

	tempPkgFile.Write(buf)
	tempPkgFile.Close()
	// real install site
	appSite := fmt.Sprintf("%s/%s", sys.PgSite(sys.PG_APP).Path,
		pc.AppName)

	// 拿到pkg.toml的信息之后，再将包解压到/pkg/app/xx
	pglib.UnpackPkg(tempPkgFile.Name(), appSite)
	// 解压到 tmp文件中，然后读取 config
	cfg, err := pglib.LoadPackageConfig(filepath.Join(appSite, pglib.PKGFILE_NAME))
	if err != nil {
		return err
	}

	// 启动成功之后，将 应用的记录，添加到 数据库中，如果已存在该app的版本，则修改版本号
	logrus.Println("add to db")

	if existApp.NodeAppId != "" {
		existApp.NodeAppVersion = pc.AppVersion
		p.appDao.Update(existApp)
	} else {
		nodeInfo, err := p.sysDao.Get()
		if err != nil {
			return fmt.Errorf("get NodeInfo,%w", err)
		}

		p.appDao.Create(&models.NodeApp{
			NodeAppId:     tool.GetUUIDUpper(),
			NodeAppType:   &models.APP_NORM,
			NodeAppName:   cfg.Name,
			NodeAppPort:   cfg.Port,
			NodeAppStatus: 0,
			NodeAppIntro:  cfg.Intro,
			NodeAppDomain: nodeInfo.GetAppDomain(strings.ToLower(cfg.Name)), // 基于此结点的根域名的，扩展子应用 app.node.pi.g
		})
	}

	return nil
}
