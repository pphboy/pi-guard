package models

import (
	"pglib/center"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	STATUS_RUNNING   = 0
	STATUS_STOP      = 1
	STATUS_ERROR     = -1
	STATUS_UN_HEALTH = -2
)

// Package models  项目
// author : pphboy@qq.com
// date : 2024-03-14 20:09
// desc : 项目

// PiProject  项目。
// 说明:
// 表名:PI_PROJECT
// group: PiProject
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.PiProject
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.PiProject
// version:2024-03-14 20:09
type PiProject struct {
	ProjectId     *int       `gorm:"column:PROJECT_ID;primaryKey;" json:"projectId"` //type:*int         comment:项目ID                               version:2024-03-14 20:09
	ProjectName   string     `gorm:"column:PROJECT_NAME;unique" json:"projectName"`  //type:string       comment:项目名                               version:2024-03-14 20:09
	ProjectStatus *int       `gorm:"column:PROJECT_STATUS" json:"projectStatus"`     //type:*int         comment:项目状态;running=1,stop=0,error=-1   version:2024-03-14 20:09
	ProjectIntro  string     `gorm:"column:PROJECT_INTRO" json:"projectIntro"`       //type:string       comment:项目介绍                             version:2024-03-14 20:09
	CreatedAt     *time.Time `gorm:"column:CREATED_AT" json:"createdAt"`             //type:*time.Time   comment:创建时间                             version:2024-03-14 20:09
	Domain        string     `gorm:"column:Domain;unique" json:"domain"`             //type:string       comment:                                     version:2024-03-15 22:58
	Port          *int       `gorm:"column:Port;unique" json:"port"`                 //type:*int         comment:                                     version:2024-03-15 22:58
	UpdatedAt     *time.Time `gorm:"column:UPDATED_AT" json:"updatedAt"`             //type:*time.Time   comment:更新时间                             version:2024-03-14 20:09
	DeletedAt     *time.Time `gorm:"column:DELETED_AT" json:"deletedAt"`             //type:*time.Time   comment:删除时间                             version:2024-03-14 20:09
}

func (p *PiProject) Msg() *center.PiProject {

	pg := &center.PiProject{
		ProjectId:     int32(*p.ProjectId),
		ProjectName:   p.ProjectName,
		ProjectStatus: int32(*p.ProjectStatus),
		ProjectIntro:  p.ProjectIntro,
		Domain:        p.Domain,
		Port:          int32(*p.Port),
	}
	if p.CreatedAt != nil {
		pg.CreatedAt = timestamppb.New(*p.CreatedAt)
	}
	if p.UpdatedAt != nil {
		pg.UpdatedAt = timestamppb.New(*p.UpdatedAt)
	}
	if p.DeletedAt != nil {
		pg.DeletedAt = timestamppb.New(*p.DeletedAt)
	}

	return pg
}

// TableName 表名:PI_PROJECT，项目。
// 说明:
func (*PiProject) TableName() string {
	return "PI_PROJECT"
}

type PiNode struct {
	NodeId     string     `gorm:"column:NODE_ID;primaryKey;" json:"nodeId"` //type:string       comment:结点编号                             version:2024-03-18 21:18
	ProjectId  *int       `gorm:"column:PROJECT_ID;" json:"projectId"`      //type:*int         comment:项目ID                               version:2024-03-18 21:18
	NodeName   string     `gorm:"column:NODE_NAME" json:"nodeName"`         //type:string       comment:结点名                               version:2024-03-18 21:18
	NodeDomain string     `gorm:"column:NODE_DOMAIN" json:"nodeDomain"`     //type:string       comment:结点域名                             version:2024-03-18 21:18
	NodeIntro  string     `gorm:"column:NODE_INTRO" json:"nodeIntro"`       //type:string       comment:结点介绍                             version:2024-03-18 21:18
	NodeStatus *int       `gorm:"column:NODE_STATUS" json:"nodeStatus"`     //type:*int         comment:结点状态;running=1,stop=0,error=-1   version:2024-03-18 21:18
	CreatedAt  *time.Time `gorm:"column:CREATED_AT" json:"createdAt"`       //type:*time.Time   comment:创建时间                             version:2024-03-18 21:18
	UpdatedAt  *time.Time `gorm:"column:UPDATED_AT" json:"updatedAt"`       //type:*time.Time   comment:更新时间                             version:2024-03-18 21:18
	DeletedAt  *time.Time `gorm:"column:DELETED_AT" json:"deletedAt"`       //type:*time.Time   comment:删除时间                             version:2024-03-18 21:18
}

// TableName 表名:PI_NODE，结点。
func (*PiNode) TableName() string {
	return "PI_NODE"
}
