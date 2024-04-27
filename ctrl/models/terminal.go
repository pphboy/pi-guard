package models

import "time"

// version:2024-03-27 21:50
type PiScript struct {
	ScriptId      string     `gorm:"column:SCRIPT_ID;primaryKey;" json:"scriptId"` //type:string       comment:脚本ID      version:2024-03-27 21:50
	ScriptName    string     `gorm:"column:SCRIPT_NAME;unique" json:"scriptName"`  //type:string       comment:脚本名称    version:2024-03-27 21:50
	ScriptContent string     `gorm:"column:SCRIPT_CONTENT" json:"scriptContent"`   //type:string       comment:脚本内容    version:2024-03-27 21:50
	CreatedAt     *time.Time `gorm:"column:CREATED_AT" json:"createdAt"`           //type:*time.Time   comment:创建时间    version:2024-03-27 21:50
	UpdatedAt     *time.Time `gorm:"column:UPDATED_AT" json:"updatedAt"`           //type:*time.Time   comment:更新时间    version:2024-03-27 21:50
	DeletedAt     *time.Time `gorm:"column:DELETED_AT" json:"deletedAt"`           //type:*time.Time   comment:删除时间    version:2024-03-27 21:50
}

// TableName 表名:PI_SCRIPT，脚本。
// 说明:
func (*PiScript) TableName() string {
	return "PI_SCRIPT"
}

// PiTerminal  终端助手。
// 说明:
// 表名:PI_TERMINAL
// group: PiTerminal
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.PiTerminal
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.PiTerminal
// version:2024-03-27 21:51
type PiTerminal struct {
	TerminalId   string     `gorm:"column:TERMINAL_ID;primaryKey;" json:"terminalId"` //type:string       comment:                                            version:2024-03-27 21:51
	TerminalType string     `gorm:"column:TERMINAL_TYPE" json:"terminalType"`         //type:string       comment:node 连接结点的, normal 是自己添加的连接    version:2024-03-27 21:51
	NodeId       string     `gorm:"column:NODE_ID" json:"nodeId"`                     //type:string       comment:                                            version:2024-03-27 21:51
	TerminalUser string     `gorm:"column:TERMINAL_USER" json:"terminalUser"`         //type:string       comment:                                            version:2024-03-27 21:51
	TerminalPswd string     `gorm:"column:TERMINAL_PSWD" json:"terminalPswd"`         //type:string       comment:                                            version:2024-03-27 21:51
	DeletedAt    *time.Time `gorm:"column:DELETED_AT" json:"deletedAt"`               //type:*time.Time   comment:删除时间                                    version:2024-03-27 21:51
	UpdatedAt    *time.Time `gorm:"column:UPDATED_AT" json:"updatedAt"`               //type:*time.Time   comment:更新时间                                    version:2024-03-27 21:51
	CreatedAt    *time.Time `gorm:"column:CREATED_AT" json:"createdAt"`               //type:*time.Time   comment:创建时间                                    version:2024-03-27 21:51
}

// TableName 表名:PI_TERMINAL，终端助手。
// 说明:
func (*PiTerminal) TableName() string {
	return "PI_TERMINAL"
}
