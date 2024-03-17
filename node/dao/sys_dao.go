package dao

import (
	"go-node/models"

	"gorm.io/gorm"
)

type SysDao interface {
	Init(ns *models.NodeSys) error
	Get() (ns *models.NodeSys, err error)
}

type SysDaoImpl struct {
	Db *gorm.DB
}

func (s *SysDaoImpl) Init(ns *models.NodeSys) error {
	return s.Db.Create(ns).Error
}

func (s *SysDaoImpl) Get() (ns *models.NodeSys, err error) {
	nodeS := &models.NodeSys{}
	err = s.Db.First(nodeS).Error
	if err != nil {
		return nodeS, nil
	}
	return
}
