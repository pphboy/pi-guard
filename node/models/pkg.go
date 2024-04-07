// Package models  结点应用
// author : p
// date : 2024-02-16 19:38
// desc : 结点应用
package models

import (
	"snproto"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	NodeAppId      string     `gorm:"column:NODE_APP_ID;primaryKey;" json:"NodeAppId"` //type:string       comment:应用ID                                         version:2024-02-16 19:38
	NodeId         string     `gorm:"column:NODE_ID" json:"NodeId"`                    //type:string       comment:结点ID                                         version:2024-02-16 19:38
	NodeAppType    *AppType   `gorm:"column:NODE_APP_TYPE" json:"NodeAppType"`         //type:*int         comment:结点应用类型;1 SYS, 2 NORMAL                   version:2024-02-16 19:38
	NodeAppName    string     `gorm:"column:NODE_APP_NAME" json:"NodeAppName"`         //type:string       comment:应用名                                         version:2024-02-16 19:38
	NodeAppIntro   string     `gorm:"column:NODE_APP_INTRO" json:"NodeAppIntro"`       //type:string       comment:应用介绍                                       version:2024-02-16 19:38
	NodeAppPort    string     `gorm:"column:NODE_APP_PORT" json:"NodeAppPort"`         //type:string       comment:应用端口                                       version:2024-02-16 19:38
	NodeAppDomain  string     `gorm:"column:NODE_APP_DOMAIN" json:"NodeAppDomain"`     //type:string       comment:应用访问地址                                   version:2024-02-16 19:38
	NodeAppStatus  int        `gorm:"column:NODE_APP_STATUS" json:"NodeAppStatus"`     //type:*int         comment:活动状态;running=1,stop=0,error=-1,unheal=-2   version:2024-02-16 19:38
	NodeAppVersion int        `gorm:"column:NODE_APP_VERSION" json:"NodeAppVersion"`   //type:*int         comment:版本                                           version:2024-02-17 16:22
	CreatedAt      *time.Time `gorm:"column:CREATED_AT" json:"CreatedAt"`              //type:*time.Time   comment:创建时间                                       version:2024-02-16 19:38
	DeletedAt      *time.Time `gorm:"column:DELETED_AT" json:"DeletedAt"`              //type:*time.Time   comment:删除时间                                       version:2024-02-16 19:38
	UpdatedAt      *time.Time `gorm:"column:UPDATED_AT" json:"UpdatedAt"`              //type:*time.Time   comment:更新时间                                       version:2024-02-16 19:38
}

func (v *NodeApp) Message() *snproto.NodeAppInfo {
	return &snproto.NodeAppInfo{
		NodeAppId:      v.NodeAppId,
		NodeId:         v.NodeId,
		NodeAppName:    v.NodeAppName,
		NodeAppType:    int32(*v.NodeAppType),
		NodeAppIntro:   v.NodeAppIntro,
		NodeAppPort:    v.NodeAppPort,
		NodeAppDomain:  v.NodeAppDomain,
		NodeAppStatus:  int32(v.NodeAppStatus),
		NodeAppVersion: int32(v.NodeAppVersion),

		CreatedAt: timestamppb.New(*v.CreatedAt),
		DeletedAt: timestamppb.New(*v.DeletedAt),
		UpdatedAt: timestamppb.New(*v.UpdatedAt),
	}
}

type AppType int

var (
	APP_SYS  AppType = 1
	APP_NORM AppType = 2
)

// TableName 表名:NODE_APP，结点应用。
// 说明:
func (*NodeApp) TableName() string {
	return "NODE_APP"
}

// PiCloudApp  应用。
// 说明:
// 表名:PI_CLOUD_APP
// group: PiCloudApp
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.PiCloudApp
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.PiCloudApp
// version:2024-02-17 09:59
type PiCloudApp struct {
	AppId      string     `gorm:"column:APP_ID;primaryKey;" json:"AppId"` //type:string       comment:应用ID                                        version:2024-02-17 09:59
	AppName    string     `gorm:"column:APP_NAME" json:"AppName"`         //type:string       comment:应用名称                                      version:2024-02-17 09:59
	AppIntro   string     `gorm:"column:APP_INTRO" json:"AppIntro"`       //type:string       comment:应用介绍                                      version:2024-02-17 09:59
	AppManual string     `gorm:"column:APP_MANUNAL" json:"AppManunal"`   //type:string       comment:应用手册                                      version:2024-02-17 09:59
	AppVersion int        `gorm:"column:APP_VERSION" json:"AppVersion"`   //type:*int         comment:应用版本                                      version:2024-02-17 16:22
	AppSite    string     `gorm:"column:APP_SITE" json:"AppSite"`         //type:string       comment:下载地址                                      version:2024-02-17 09:59
	AppHistory string     `gorm:"column:APP_HISTORY" json:"AppHistory"`   //type:string       comment:应用历史;存储过往版本名与时间，以及下载方式   version:2024-02-17 09:59
	CreatedAt  *time.Time `gorm:"column:CREATED_AT" json:"CreatedAt"`     //type:*time.Time   comment:创建时间                                      version:2024-02-17 09:59
	UpdatedAt  *time.Time `gorm:"column:UPDATED_AT" json:"UpdatedAt"`     //type:*time.Time   comment:更新时间                                      version:2024-02-17 09:59
	DeletedAt  *time.Time `gorm:"column:DELETED_AT" json:"DeletedAt"`     //type:*time.Time   comment:删除时间                                      version:2024-02-17 09:59
}

// TableName 表名:PI_CLOUD_APP，应用。
// 说明:
func (*PiCloudApp) TableName() string {
	return "PI_CLOUD_APP"
}
