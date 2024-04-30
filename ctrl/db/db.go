package db

import (
	"fmt"
	cm "go-ctrl/models"
	"go-node/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	cdb *gorm.DB
)

func Init(dataPath string, name string) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/%s.db", dataPath, name)), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&cm.PiNode{},
		&cm.PiTerminal{},
		&cm.PiScript{},
		&cm.PiManager{},
		&cm.PiProject{},
		&models.PiCloudApp{}); err != nil {
		panic(err)
	}

	cdb = db

}

func DB() *gorm.DB {
	if cdb == nil {
		panic("db is not initialized!!!")
	}
	return cdb
}
