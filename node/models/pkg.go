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
	NodeAppId      string     `gorm:"column:NODE_APP_ID;primaryKey;" json:"appId"` //type:string       comment:应用ID                                         version:2024-02-16 19:38
	NodeId         string     `gorm:"column:NODE_ID" json:"nodeId"`                //type:string       comment:结点ID                                         version:2024-02-16 19:38
	NodeAppType    *AppType   `gorm:"column:NODE_APP_TYPE" json:"appType"`         //type:*int         comment:结点应用类型;1 SYS, 2 NORMAL                   version:2024-02-16 19:38
	NodeAppName    string     `gorm:"column:NODE_APP_NAME" json:"appName"`         //type:string       comment:应用名                                         version:2024-02-16 19:38
	NodeAppIntro   string     `gorm:"column:NODE_APP_INTRO" json:"appIntro"`       //type:string       comment:应用介绍                                       version:2024-02-16 19:38
	NodeAppPort    string     `gorm:"column:NODE_APP_PORT" json:"appPort"`         //type:string       comment:应用端口                                       version:2024-02-16 19:38
	NodeAppDomain  string     `gorm:"column:NODE_APP_DOMAIN" json:"appDomain"`     //type:string       comment:应用访问地址                                   version:2024-02-16 19:38
	NodeAppStatus  int        `gorm:"column:NODE_APP_STATUS" json:"appStatus"`     //type:*int         comment:活动状态;running=1,stop=0,error=-1,unheal=-2   version:2024-02-16 19:38
	NodeAppVersion int        `gorm:"column:NODE_APP_VERSION" json:"appVersion"`   //type:*int         comment:版本                                           version:2024-02-17 16:22
	CreatedAt      *time.Time `gorm:"column:CREATED_AT" json:"CreatedAt"`          //type:*time.Time   comment:创建时间                                       version:2024-02-16 19:38
	DeletedAt      *time.Time `gorm:"column:DELETED_AT" json:"DeletedAt"`          //type:*time.Time   comment:删除时间                                       version:2024-02-16 19:38
	UpdatedAt      *time.Time `gorm:"column:UPDATED_AT" json:"UpdatedAt"`          //type:*time.Time   comment:更新时间                                       version:2024-02-16 19:38
}

func (v *NodeApp) Message() *snproto.NodeAppInfo {

	nai := &snproto.NodeAppInfo{
		NodeAppId:   v.NodeAppId,
		NodeId:      v.NodeId,
		NodeAppName: v.NodeAppName,
		// NodeAppType:    int32(*v.NodeAppType),
		NodeAppIntro:   v.NodeAppIntro,
		NodeAppPort:    v.NodeAppPort,
		NodeAppDomain:  v.NodeAppDomain,
		NodeAppStatus:  int32(v.NodeAppStatus),
		NodeAppVersion: int32(v.NodeAppVersion),
	}

	if v.NodeAppType != nil {
		nai.NodeAppType = int32(*v.NodeAppType)
	}
	if v.CreatedAt != nil {
		nai.CreatedAt = timestamppb.New(*v.CreatedAt)
	}
	if v.UpdatedAt != nil {
		nai.UpdatedAt = timestamppb.New(*v.UpdatedAt)
	}
	if v.DeletedAt != nil {
		nai.DeletedAt = timestamppb.New(*v.DeletedAt)
	}

	return nai
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
	AppId      string     `gorm:"column:APP_ID;primaryKey;" json:"appId"` //type:string       comment:应用ID                                        version:2024-02-17 09:59
	AppName    string     `gorm:"column:APP_NAME;unique;" json:"appName"` //type:string       comment:应用名称                                      version:2024-02-17 09:59
	AppIntro   string     `gorm:"column:APP_INTRO" json:"appIntro"`       //type:string       comment:应用介绍                                      version:2024-02-17 09:59
	AppManual  string     `gorm:"column:APP_MANUAL" json:"appManual"`     //type:string       comment:应用手册                                      version:2024-02-17 09:59
	AppVersion int        `gorm:"column:APP_VERSION" json:"appVersion"`   //type:*int         comment:应用版本                                      version:2024-02-17 16:22
	AppSite    string     `gorm:"column:APP_SITE;unique;" json:"appSite"` //type:string       comment:下载地址                                      version:2024-02-17 09:59
	AppHistory string     `gorm:"column:APP_HISTORY" json:"appHistory"`   //type:string       comment:应用历史;存储过往版本名与时间，以及下载方式   version:2024-02-17 09:59
	CreatedAt  *time.Time `gorm:"column:CREATED_AT" json:"createdAt"`     //type:*time.Time   comment:创建时间                                      version:2024-02-17 09:59
	UpdatedAt  *time.Time `gorm:"column:UPDATED_AT" json:"updatedAt"`     //type:*time.Time   comment:更新时间                                      version:2024-02-17 09:59
	DeletedAt  *time.Time `gorm:"column:DELETED_AT" json:"deletedAt"`     //type:*time.Time   comment:删除时间                                      version:2024-02-17 09:59
}

// TableName 表名:PI_CLOUD_APP，应用。
// 说明:
func (*PiCloudApp) TableName() string {
	return "PI_CLOUD_APP"
}

// 将grpc的消息转成models的
func (*PiCloudApp) ResolveGrpcMsg(gm *snproto.PiCloudApp) *PiCloudApp {
	ct := gm.CreatedAt.AsTime()
	ut := gm.UpdatedAt.AsTime()
	dt := gm.DeletedAt.AsTime()
	return &PiCloudApp{
		AppId:      gm.AppId,
		AppName:    gm.AppName,
		AppIntro:   gm.AppIntro,
		AppManual:  gm.AppManual,
		AppVersion: int(gm.AppVersion),
		AppSite:    gm.AppSite,
		AppHistory: gm.AppHistory,
		CreatedAt:  &ct,
		UpdatedAt:  &ut,
		DeletedAt:  &dt,
	}
}

func (gm *PiCloudApp) GrpcMsg() *snproto.PiCloudApp {

	pca := &snproto.PiCloudApp{
		AppId:      gm.AppId,
		AppName:    gm.AppName,
		AppIntro:   gm.AppIntro,
		AppManual:  gm.AppManual,
		AppVersion: int32(gm.AppVersion),
		AppSite:    gm.AppSite,
		AppHistory: gm.AppHistory,
	}

	if gm.CreatedAt != nil {
		pca.CreatedAt = timestamppb.New(*gm.CreatedAt)
	}

	if gm.UpdatedAt != nil {
		pca.UpdatedAt = timestamppb.New(*gm.UpdatedAt)
	}

	if gm.DeletedAt != nil {
		pca.DeletedAt = timestamppb.New(*gm.DeletedAt)
	}

	return pca
}
