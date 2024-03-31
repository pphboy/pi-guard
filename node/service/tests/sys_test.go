package service

import (
	"go-node/service"
	"testing"
)

func TestSysServiceInstall(t *testing.T) {
	i := service.NewIniter("../../fs_root")
	sysService := service.NewSysService(service.BaseService{
		DB: i.GetDB(),
	})

	if err := sysService.Install("ndev"); err != nil {
		t.Fatal(err)
	}

	info, err := sysService.GetSysInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)

}
