// Package models  结点应用
// author : p
// date : 2024-02-16 19:38
// desc : 结点应用
package models

import "time"

// NodeApp  结点应用。
// 说明:
// 表名:NODE_APP
// group: NodeApp
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.NodeApp
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.NodeApp
// version:2024-02-16 19:38
type NodeApp struct {
	NodeAppId     string     `gorm:"column:primaryKey;NODE_APP_ID" json:"NodeAppId"` //type:string       comment:应用ID                                         version:2024-02-16 19:38
	NodeId        string     `gorm:"column:NODE_ID" json:"NodeId"`                   //type:string       comment:结点ID                                         version:2024-02-16 19:38
	NodeAppType   *int       `gorm:"column:NODE_APP_TYPE" json:"NodeAppType"`        //type:*int         comment:结点应用类型;1 SYS, 2 NORMAL                   version:2024-02-16 19:38
	NodeAppName   string     `gorm:"column:NODE_APP_NAME" json:"NodeAppName"`        //type:string       comment:应用名                                         version:2024-02-16 19:38
	NodeAppIntro  string     `gorm:"column:NODE_APP_INTRO" json:"NodeAppIntro"`      //type:string       comment:应用介绍                                       version:2024-02-16 19:38
	NodeAppPort   string     `gorm:"column:NODE_APP_PORT" json:"NodeAppPort"`        //type:string       comment:应用端口                                       version:2024-02-16 19:38
	NodeAppDomain string     `gorm:"column:NODE_APP_DOMAIN" json:"NodeAppDomain"`    //type:string       comment:应用访问地址                                   version:2024-02-16 19:38
	NodeAppStatus *int       `gorm:"column:NODE_APP_STATUS" json:"NodeAppStatus"`    //type:*int         comment:活动状态;running=1,stop=0,error=-1,unheal=-2   version:2024-02-16 19:38
	CreatedAt     *time.Time `gorm:"column:CREATED_AT" json:"CreatedAt"`             //type:*time.Time   comment:创建时间                                       version:2024-02-16 19:38
	DeletedAt     *time.Time `gorm:"column:DELETED_AT" json:"DeletedAt"`             //type:*time.Time   comment:删除时间                                       version:2024-02-16 19:38
	UpdatedAt     *time.Time `gorm:"column:UPDATED_AT" json:"UpdatedAt"`             //type:*time.Time   comment:更新时间                                       version:2024-02-16 19:38
}

// TableName 表名:NODE_APP，结点应用。
// 说明:
func (*NodeApp) TableName() string {
	return "NODE_APP"
}
