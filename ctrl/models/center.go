package models

import (
	"pglib/center"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
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
	ProjectId     *int       `gorm:"column:primaryKey;PROJECT_ID" json:"ProjectId"` //type:*int         comment:项目ID                               version:2024-03-14 20:09
	ProjectName   string     `gorm:"column:PROJECT_NAME" json:"ProjectName"`        //type:string       comment:项目名                               version:2024-03-14 20:09
	ProjectStatus *int       `gorm:"column:PROJECT_STATUS" json:"ProjectStatus"`    //type:*int         comment:项目状态;running=1,stop=0,error=-1   version:2024-03-14 20:09
	ProjectIntro  string     `gorm:"column:PROJECT_INTRO" json:"ProjectIntro"`      //type:string       comment:项目介绍                             version:2024-03-14 20:09
	CreatedAt     *time.Time `gorm:"column:CREATED_AT" json:"CreatedAt"`            //type:*time.Time   comment:创建时间                             version:2024-03-14 20:09
	Domain        string     `gorm:"column:Domain" json:"Domain"`                   //type:string       comment:                                     version:2024-03-15 22:58
	Port          *int       `gorm:"column:Port" json:"Port"`                       //type:*int         comment:                                     version:2024-03-15 22:58
	UpdatedAt     *time.Time `gorm:"column:UPDATED_AT" json:"UpdatedAt"`            //type:*time.Time   comment:更新时间                             version:2024-03-14 20:09
	DeletedAt     *time.Time `gorm:"column:DELETED_AT" json:"DeletedAt"`            //type:*time.Time   comment:删除时间                             version:2024-03-14 20:09
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
