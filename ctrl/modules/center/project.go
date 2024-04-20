package centers

import (
	"go-ctrl/db"
	"go-ctrl/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProjectManager struct {
	db      *gorm.DB
	centers *Manager
}

func NewProjectManager(m *Manager) *ProjectManager {
	pm := &ProjectManager{
		db:      db.DB(),
		centers: m,
	}
	// 启动结点
	go pm.bootCenters()

	return pm
}

func (p *ProjectManager) bootCenters() {
	plist, err := p.List()
	if err != nil {
		logrus.Fatal("booting project list", err)
	}
	p.centers.BootCenters(plist)

	for _, v := range plist {
		if err := p.SetStatus(v, models.STATUS_RUNNING); err != nil {
			logrus.Errorf("%+v, failed to start", v)
			continue
		}
	}
}

// 完成新建项目，然后就是添加接口，查看结点，然后就是连接ssh，传输文件等
func (p *ProjectManager) Create(pj *models.PiProject) error {
	if err := p.db.Create(pj).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProjectManager) Delete(pj *models.PiProject) error {
	if err := p.db.Delete(pj, "project_id = ?", pj.ProjectId).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProjectManager) List() (ps []*models.PiProject, err error) {
	if err := p.db.Find(&ps).Error; err != nil {
		return nil, err
	}

	for k, v := range ps {
		if v.ProjectStatus == nil {
			ps[k].ProjectStatus = &(models.STATUS_STOP)
		}
	}
	return
}

func (p *ProjectManager) SetStatus(ps *models.PiProject, status int) error {
	ps.ProjectStatus = &status
	if err := p.db.Where("project_id = ?", ps.ProjectId).Updates(ps).Error; err != nil {
		return err
	}
	return nil
}
