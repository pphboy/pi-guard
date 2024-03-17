// Package models  结点系统
// author : p
// date : 2024-02-16 19:36
// desc : 结点系统
package models

import (
	"fmt"
	"time"
)

// NodeSys  结点系统。
// 说明:
// 表名:NODE_SYS
// group: NodeSys
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.NodeSys
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.NodeSys
// version:2024-02-16 19:36
type NodeSys struct {
	NodeId     string     `gorm:"column:NODE_ID;primaryKey;" json:"NodeId"` //type:string       comment:                             version:2024-02-16 19:36
	NodeName   string     `gorm:"column:NODE_NAME" json:"NodeName"`         //type:string       comment:                             version:2024-02-16 19:36
	NodeStatus *int       `gorm:"column:NODE_STATUS" json:"NodeStatus"`     //type:*int         comment:running=1,stop=0,error=-1    version:2024-02-16 19:36
	NodeDomain string     `gorm:"column:NODE_DOMAIN" json:"NodeDomain"`     //type:string       comment:                             version:2024-02-17 16:32
	CreatedAt  *time.Time `gorm:"column:CREATED_AT" json:"CreatedAt"`       //type:*time.Time   comment:创建时间                     version:2024-02-16 19:36
	UpdatedAt  *time.Time `gorm:"column:UPDATED_AT" json:"UpdatedAt"`       //type:*time.Time   comment:更新时间                     version:2024-02-16 19:36
	DeletedAt  *time.Time `gorm:"column:DELETED_AT" json:"DeletedAt"`       //type:*time.Time   comment:删除时间                     version:2024-02-16 19:36
}

func (n *NodeSys) GetAppDomain(app string) string {
	return fmt.Sprintf("%s.%s", app, n.NodeDomain)
}

// TableName 表名:NODE_SYS，结点系统。
// 说明:
func (*NodeSys) TableName() string {
	return "NODE_SYS"
}

// NodeSetting  结点设置。
// 说明:
// 表名:NODE_SETTING
// group: NodeSetting
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.NodeSetting
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.NodeSetting
// version:2024-02-16 19:37
type NodeSetting struct {
	SettingId      string     `gorm:"column:SETTING_ID;primaryKey;" json:"SettingId"` //type:string       comment:                          version:2024-02-16 19:37
	SettingKey     string     `gorm:"column:SETTING_KEY" json:"SettingKey"`           //type:string       comment:                          version:2024-02-16 19:37
	SettingContent string     `gorm:"column:SETTING_CONTENT" json:"SettingContent"`   //type:string       comment:设置的内容，可存json等    version:2024-02-16 19:37
	CreatedAt      *time.Time `gorm:"column:CREATED_AT" json:"CreatedAt"`             //type:*time.Time   comment:创建时间                  version:2024-02-16 19:37
	UpdatedAt      *time.Time `gorm:"column:UPDATED_AT" json:"UpdatedAt"`             //type:*time.Time   comment:更新时间                  version:2024-02-16 19:37
	DeletedAt      *time.Time `gorm:"column:DELETED_AT" json:"DeletedAt"`             //type:*time.Time   comment:删除时间                  version:2024-02-16 19:37
}

// TableName 表名:NODE_SETTING，结点设置。
// 说明:
func (*NodeSetting) TableName() string {
	return "NODE_SETTING"
}
