package modules

import (
	"go-node/models"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

type ServiceManager interface {
	GetServices() []RunAppService
	GetServiceByApp(*models.NodeApp) RunAppService
}

type RunAppService interface {
	Start() error
	Stop() error
	Restart() error
}

type RunnerApp struct {
	Cmd     *exec.Cmd
	App     *models.NodeApp
	reMutex sync.Mutex
}

func (r *RunnerApp) Start() error {
	go func() {
		if err := r.Cmd.Run(); err != nil {
			logrus.Errorf("App:%s, Cmd:%s, Err: %v ", r.App.NodeAppName, r.Cmd.String(), err)
		}
	}()
	return nil
}

func (r *RunnerApp) Stop() error {
	if err := r.Cmd.Process.Kill(); err != nil {
		return err
	}
	if err := r.Cmd.Process.Release(); err != nil {
		return err
	}
	return nil
}

func (r *RunnerApp) Restart() error {
	r.reMutex.Lock()
	defer r.reMutex.Unlock()
	r.Cmd = exec.Command(r.Cmd.String())
	return r.Start()
}
