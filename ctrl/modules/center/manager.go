package centers

import (
	"go-ctrl/db"
	"go-ctrl/models"
	"pglib/cdns"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Manager struct {
	list   []Center
	group  *gin.RouterGroup
	db     *gorm.DB
	dnsM   cdns.CdnsManager
	netseg string
}

func NewManager(group *gin.RouterGroup, dnsM cdns.CdnsManager, netseg string) *Manager {
	return &Manager{
		group:  group,
		db:     db.DB(),
		dnsM:   dnsM,
		netseg: netseg,
	}
}

func (m *Manager) AddCenter(p *models.PiProject) {
	m.list = append(m.list, NewCenter(m.db, p, m.group, m.dnsM, m.netseg))
}

func (m *Manager) BootCenters(p []*models.PiProject) {
	for _, v := range p {
		m.list = append(m.list, NewCenter(m.db, v, m.group, m.dnsM, m.netseg))
	}
}
