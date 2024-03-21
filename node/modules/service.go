package modules

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
	GetServiceByApp(*models.NodeApp) RunAppService
	AddServices(RunAppService)
}

type RunAppService interface {
	Start() error
	Stop() error
	Restart() error
	Equal(*models.NodeApp) bool
}

type AppDirector struct {
	services []RunAppService
	mutex    sync.Mutex
}

func (a *AppDirector) GetServices() []RunAppService {
	return nil
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

func (a *AppDirector) AddServices(app RunAppService) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.services = append(a.services, app)
}

func NewRunnerApp(cmd *exec.Cmd, app *models.NodeApp) RunAppService {
	return &RunnerApp{
		cmd: cmd,
		app: app,
	}
}

type RunnerApp struct {
	cmd     *exec.Cmd
	app     *models.NodeApp
	reMutex sync.Mutex
}

func (r *RunnerApp) Start() error {
	go func() {
		if err := r.cmd.Run(); err != nil {
			logrus.Errorf("App:%s, Cmd:%s, Err: %v ", r.app.NodeAppName, r.cmd.String(), err)
		}
	}()
	return nil
}

func (r *RunnerApp) Stop() error {
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
