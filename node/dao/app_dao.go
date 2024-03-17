package dao

import (
	"go-node/models"

	"gorm.io/gorm"
)

type AppDao interface {
	Create(app *models.NodeApp) error
	Update(app *models.NodeApp) error
	Delete(app *models.NodeApp) error
	GetAll() ([]*models.NodeApp, error)
	GetByName(name string) (*models.NodeApp, error)
}

type AppDaoImpl struct {
	Db *gorm.DB
}

func (a *AppDaoImpl) Create(app *models.NodeApp) error {
	return a.Db.Debug().FirstOrCreate(app).Error
}

func (a *AppDaoImpl) GetByName(name string) (app *models.NodeApp, err error) {
	app = &models.NodeApp{}
	err = a.Db.Model(app).Where("NODE_APP_NAME = ?", name).First(app).Error
	return
}

func (a *AppDaoImpl) Update(app *models.NodeApp) error {
	return a.Db.Debug().Where("NODE_APP_ID = ?", app.NodeAppId).Updates(app).Error
}

func (a *AppDaoImpl) Delete(app *models.NodeApp) error {
	return a.Db.Debug().Delete(app, "NODE_APP_ID = ?", app.NodeAppId).Error
}

func (a *AppDaoImpl) GetAll() (list []*models.NodeApp, err error) {
	err = a.Db.Find(&list).Error
	return
}
