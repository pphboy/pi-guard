package scripter

import (
	"go-ctrl/db"
	"go-ctrl/models"
	"go-node/tool"

	"gorm.io/gorm"
)

type ScripterManager struct {
	db *gorm.DB
}

func NewScripterManager() *ScripterManager {
	return &ScripterManager{
		db: db.DB(),
	}
}

func (s *ScripterManager) Create(t *models.PiScript) error {
	t.ScriptId = tool.GetUUIDUpper()
	if err := s.db.FirstOrCreate(t).Error; err != nil {
		return err
	}
	return nil
}

func (s *ScripterManager) Delete(id string) error {
	if err := s.db.Unscoped().Where("script_id = ?", id).Delete(&models.PiScript{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *ScripterManager) List() (ts []*models.PiScript, err error) {
	if err := s.db.Model(&models.PiScript{}).Find(&ts).Error; err != nil {
		return nil, err
	}
	return
}
