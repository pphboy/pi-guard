package service

import (
	"fmt"
	"go-node/models"
	gsys "go-node/sys"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Initer interface {
	InitBase()
	GetDB() *gorm.DB
}

func NewIniter() Initer {
	return &InitService{}
}

type InitService struct {
	db *gorm.DB
}

// 初始化基础
func (i *InitService) InitBase() {
	i.initializeDefaultPath()
}

// 初始化路径
func (i *InitService) initializeDefaultPath() {
	siter := NewSiter(gsys.GetPgSites())
	siter.InitNodeSysPath()
}

// 加载应用服务
// func (i *InitService) loadAppService() {

// }

func (i *InitService) GetDB() *gorm.DB {
	if i.db == nil {
		// gorm.Open()
		dbPath := fmt.Sprintf("%s/%s", gsys.PgSite(gsys.PG_DB).Path, "node.db")
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			logrus.Fatal("initialize failed ,", dbPath, err)
		}

		db.AutoMigrate(&models.NodeApp{}, &models.NodeSys{}, &models.NodeSetting{})
		i.db = db
	}

	return i.db
}
