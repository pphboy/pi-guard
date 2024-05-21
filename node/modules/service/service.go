package service

import (
	"fmt"
	"go-node/models"
	"os/exec"
	"pglib/cdns"
	"pi_dns/server"
	"strings"

	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrNodeAppNotInService = fmt.Errorf("node app not in service")
)

const (
	EVENT_ADDAPP = "event_service_add_app"
)

type ServiceManager interface {
	GetServices() []RunAppService
	GetServiceByApp(*models.NodeApp) (RunAppService, error)
	AddService(RunAppService)
	GetRunning() []RunAppService
	GetClosed() []RunAppService
}

type RunAppService interface {
	Start() error
	Close() error
	Restart() error
	Equal(*models.NodeApp) bool
	IsRunning() bool
}

func NewAppDirector() ServiceManager {
	a := &AppDirector{}
	return a
}

type AppDirector struct {
	services []RunAppService
	mutex    sync.Mutex
}

func (a *AppDirector) GetRunning() []RunAppService {
	var appList []RunAppService
	for _, v := range a.services {
		if v.IsRunning() {
			appList = append(appList, v)
		}
	}
	return appList
}

func (a *AppDirector) GetClosed() []RunAppService {
	var appList []RunAppService
	for _, v := range a.services {
		if !v.IsRunning() {
			appList = append(appList, v)
		}
	}
	return appList
}

func (a *AppDirector) GetServices() []RunAppService {
	return a.services
}

func (a *AppDirector) GetServiceByApp(e *models.NodeApp) (RunAppService, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, v := range a.services {
		if v.Equal(e) {
			return v, nil
		}
	}

	return nil, ErrNodeAppNotInService
}

func (a *AppDirector) AddService(app RunAppService) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.services = append(a.services, app)

	go app.Start()
}

// cmd 运行命令，app 为 app的安装信息
func NewRunnerApp(cmd *exec.Cmd, app *models.NodeApp, db *gorm.DB, dnsM cdns.CdnsManager, nodeIp string) RunAppService {
	return &RunnerApp{
		cmd:    cmd,
		app:    app,
		db:     db,
		dnsM:   dnsM,
		nodeIp: nodeIp,
	}
}

type RunnerApp struct {
	cmd       *exec.Cmd
	nodeIp    string
	db        *gorm.DB
	app       *models.NodeApp
	reMutex   sync.Mutex
	runStatus bool
	dnsM      cdns.CdnsManager
}

func (r *RunnerApp) Start() error {
	logrus.Print("booting ", r.app.NodeAppDomain)
	r.runStatus = true
	r.cmd = exec.Command(r.cmd.String())

	go func() {
		// bind domain to ip
		if err := r.dnsM.AddHosts(server.Host{
			Domain: r.app.NodeAppDomain,
			Ips:    []string{r.nodeIp},
		}); err != nil {
			logrus.Errorf("failed binding dns %s to %s", r.app.NodeAppDomain, r.nodeIp)
		}

		if err := r.cmd.Run(); err != nil {
			logrus.Errorf("App:%s, Cmd:%s, Err: %v ", r.app.NodeAppName, r.cmd.String(), err)

			logrus.Warnf("SyscalError %+v", strings.Compare(err.Error(), "signal: killed"))

			if strings.Compare(err.Error(), "signal: killed") == 0 {
				r.runStatus = false
				// 即使 变成 error，也可以尝试去启动
				r.setStatus(models.STATUS_STOP)
			} else {
				r.runStatus = false
				// 即使 变成 error，也可以尝试去启动
				r.setStatus(models.STATUS_ERROR)
			}

			if err := r.dnsM.DelHosts(server.Host{
				Domain: r.app.NodeAppDomain,
				Ips:    []string{r.nodeIp},
			}); err != nil {
				logrus.Errorf("failed remove dns binded %s: %s", r.app.NodeAppDomain, r.nodeIp)
			}
		}
	}()

	r.setStatus(models.STATUS_RUNNING)
	return nil
}

func (r *RunnerApp) Close() error {
	defer func() {
		r.reMutex.Lock()
		defer r.reMutex.Unlock()

		r.runStatus = false
	}()
	if err := r.cmd.Process.Kill(); err != nil {
		return err
	}
	if err := r.cmd.Process.Release(); err != nil {
		return err
	}

	r.setStatus(models.STATUS_STOP)
	return nil
}

func (r *RunnerApp) Restart() error {
	r.reMutex.Lock()
	defer r.reMutex.Unlock()
	r.cmd = exec.Command(r.cmd.String())
	return r.Start()
}

func (r *RunnerApp) Equal(e *models.NodeApp) bool {
	return r.app.NodeAppId == e.NodeAppId
}

func (r *RunnerApp) IsRunning() bool {
	r.reMutex.Lock()
	defer r.reMutex.Unlock()

	return r.runStatus
}

func (r *RunnerApp) setStatus(status int) error {
	r.app.NodeAppStatus = status
	if err := r.db.Where("NODE_APP_ID = ?", r.app.NodeAppId).Updates(r.app).Error; err != nil {
		return err
	}
	return nil
}
