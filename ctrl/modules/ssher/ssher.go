package ssher

import (
	"errors"
	"go-ctrl/db"
	"go-ctrl/models"
	"go-node/tool"
	"strings"

	"gorm.io/gorm"
)

type SsherManager struct {
	db *gorm.DB
}

func NewSsherManager() *SsherManager {
	return &SsherManager{
		db: db.DB(),
	}
}

func (s *SsherManager) Create(t *models.PiTerminal) error {
	if strings.Compare(t.NodeId, "") == 0 {
		return errors.New("get nil node id")
	}
	t.TerminalId = tool.GetUUIDUpper()
	if err := s.db.FirstOrCreate(t).Error; err != nil {
		return err
	}
	return nil
}

func (s *SsherManager) Delete(id string) error {
	if err := s.db.Unscoped().Where("terminal_id = ?", id).Delete(&models.PiTerminal{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *SsherManager) List() (ts []*models.PiTerminal, err error) {
	if err := s.db.Model(&models.PiTerminal{}).Find(&ts).Error; err != nil {
		return nil, err
	}
	return
}
