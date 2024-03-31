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
	initBase(string)
	GetDB() *gorm.DB
}

func NewIniter(base string) Initer {
	i := &InitService{}
	i.initBase(base)
	return i
}

type InitService struct {
	db *gorm.DB
}

// 初始化基础
func (i *InitService) initBase(base string) {
	logrus.Print("initer base")
	i.initializeDefaultPath(base)
}

// 初始化路径
func (i *InitService) initializeDefaultPath(base string) {
	siter := NewSiter(gsys.GetPgSites(base))
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
