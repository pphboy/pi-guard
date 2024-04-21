package appm

import (
	"errors"
	"fmt"
	"go-ctrl/db"
	"go-node/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppStore struct {
	db *gorm.DB
}

func NewAppStore() *AppStore {
	return &AppStore{
		db: db.DB(),
	}
}

func (a *AppStore) Create(p *models.PiCloudApp) error {
	p.AppId = uuid.NewString()
	if err := a.db.FirstOrCreate(p).Error; err != nil {
		return err
	}
	return nil
}

func (a *AppStore) Update(p *models.PiCloudApp) error {
	old, err := a.Get(p.AppId)
	if err != nil {
		return err
	}

	if p.AppVersion <= old.AppVersion {
		return errors.New("error version , neo version less equal than old")
	}

	// 更新历史
	p.AppHistory = fmt.Sprintf("%s\n%s",
		fmt.Sprintf("name:%s, version:%d ,site:%s, time:%v", old.AppName, old.AppVersion, old.AppSite, old.UpdatedAt),
		old.AppHistory)

	if err := a.db.Where("app_id = ?", p.AppId).Updates(p).Error; err != nil {
		return err
	}

	return nil
}

func (a *AppStore) Delete(pid string) error {
	if err := a.db.Where("app_id = ?", pid).Delete(&models.PiCloudApp{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *AppStore) List() (ps []*models.PiCloudApp, err error) {
	if err := a.db.Find(&ps).Error; err != nil {
		return nil, err
	}
	return
}

func (a *AppStore) Get(pid string) (p *models.PiCloudApp, err error) {
	p = &models.PiCloudApp{}
	if err := a.db.Where("app_id = ?", pid).First(p).Error; err != nil {
		return nil, err
	}
	return
}
