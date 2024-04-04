package service

import (
	"fmt"
	"go-node/models"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	ErrNodeAppNotInService = fmt.Errorf("node app not in service")
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
	return &AppDirector{}
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
func NewRunnerApp(cmd *exec.Cmd, app *models.NodeApp) RunAppService {
	return &RunnerApp{
		cmd: cmd,
		app: app,
	}
}

type RunnerApp struct {
	cmd       *exec.Cmd
	app       *models.NodeApp
	reMutex   sync.Mutex
	runStatus bool
}

func (r *RunnerApp) Start() error {
	r.runStatus = true

	go func() {
		if err := r.cmd.Run(); err != nil {
			logrus.Errorf("App:%s, Cmd:%s, Err: %v ", r.app.NodeAppName, r.cmd.String(), err)
			r.runStatus = false
		}
	}()
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
