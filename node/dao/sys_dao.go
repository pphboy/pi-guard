package dao

import (
	"go-node/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SysDao interface {
	Init(ns *models.NodeSys) error
	Get() (ns *models.NodeSys, err error)
	IsInstalled() (bool, error)
}

func NewSysDao(db *gorm.DB) SysDao {
	return &SysDaoImpl{
		Db: db,
	}
}

type SysDaoImpl struct {
	Db *gorm.DB
}

func (s *SysDaoImpl) Init(ns *models.NodeSys) error {
	ined, err := s.IsInstalled()
	if err != nil {
		return err
	}
	if ined {
		logrus.Print("current node is installed")
		return nil
	} else {
		return s.Db.Create(ns).Error
	}
}

func (s *SysDaoImpl) Get() (ns *models.NodeSys, err error) {
	nodeS := &models.NodeSys{}
	err = s.Db.Model(nodeS).First(nodeS).Error
	if err != nil {
		return nil, err
	}
	return nodeS, nil
}

func (s *SysDaoImpl) IsInstalled() (bool, error) {
	var count int64
	if err := s.Db.Model(&models.NodeSys{}).Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
