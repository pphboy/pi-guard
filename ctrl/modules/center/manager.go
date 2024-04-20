package centers

import (
	"go-ctrl/db"
	"go-ctrl/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Manager struct {
	list  []Center
	group *gin.RouterGroup
	db    *gorm.DB
}

func NewManager(group *gin.RouterGroup) *Manager {
	return &Manager{
		group: group,
		db:    db.DB(),
	}
}

func (m *Manager) AddCenter(p *models.PiProject) {
	m.list = append(m.list, NewCenter(m.db, p, m.group))
}

func (m *Manager) BootCenters(p []*models.PiProject) {
	for _, v := range p {
		m.list = append(m.list, NewCenter(m.db, v, m.group))
	}
}
