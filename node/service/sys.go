package service

import (
	"errors"
	"fmt"
	"go-node/dao"
	"go-node/models"
	"go-node/sys"
	"go-node/tool"

	"github.com/sirupsen/logrus"
)

var (
	ErrSysNodeInitiliazed = errors.New("node is initialized")
)

type SysService interface {
	Install(name string) error
	GetSysInfo() (*models.NodeSys, error)
}

func NewSysService(bs BaseService) SysService {
	return &SysServiceImpl{
		sysDao: dao.NewSysDao(bs.DB),
	}
}

type SysServiceImpl struct {
	sysDao dao.SysDao
}

func (s *SysServiceImpl) Install(nodeName string) error {
	sys := &models.NodeSys{
		NodeId:     tool.GetUUIDUpper(),
		NodeStatus: 0,
		NodeName:   nodeName,
		NodeDomain: fmt.Sprintf("%s.%s", nodeName, sys.ROOT_DOMAIN),
	}

	logrus.Printf("install %+v \n", sys)

	if err := s.sysDao.Init(sys); err != nil {
		return err
	}

	// 只会执行一次
	return nil
}

func (s *SysServiceImpl) GetSysInfo() (*models.NodeSys, error) {
	return s.sysDao.Get()
}
