package models

import "time"

type ManagerType int

const (
	USER_SUPER  ManagerType = 1
	USER_ADMIN  ManagerType = 2
	USER_NORMAL ManagerType = 3
)

// version:2024-03-29 20:36
type PiManager struct {
	ManagerId       string      `gorm:"column:MANAGER_ID;primaryKey;" json:"managerId"`        //type:string       comment:管理员ID                                    version:2024-03-29 20:36
	ManagerUsername string      `gorm:"column:MANAGER_USERNAME;unique" json:"managerUsername"` //type:string       comment:管理员用户名                                version:2024-03-29 20:36
	ManagerSolt     string      `gorm:"column:MANAGER_SOLT" json:"managerSolt"`                //type:*int         comment:加密盐;running=1,stop=0,error=-1            version:2024-03-29 20:36
	ManagerPswd     string      `gorm:"column:MANAGER_PSWD" json:"managerPswd"`                //type:string       comment:密码                                        version:2024-03-29 20:36
	ManagerType     ManagerType `gorm:"column:MANAGER_TYPE" json:"managerType"`                //type:string       comment:管理员类型;super = 0, admin = 1, user = 2   version:2024-03-29 20:36
	CreatedAt       *time.Time  `gorm:"column:CREATED_AT" json:"createdAt"`                    //type:*time.Time   comment:创建时间                                    version:2024-03-29 20:36
	UpdatedAt       *time.Time  `gorm:"column:UPDATED_AT" json:"updatedAt"`                    //type:*time.Time   comment:更新时间                                    version:2024-03-29 20:36
	DeletedAt       *time.Time  `gorm:"column:DELETED_AT" json:"deletedAt"`                    //type:*time.Time   comment:删除时间                                    version:2024-03-29 20:36
}

// TableName 表名:PI_MANAGER，管理员。
// 说明:
func (*PiManager) TableName() string {
	return "PI_MANAGER"
}
