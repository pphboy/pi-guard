package service

import (
	gsys "go-node/sys"
)

type Initer interface {
}

type InitService struct{}

// 初始化基础
func (i *InitService) InitBase() {
}

// 初始化路径
func (i *InitService) initializeDefaultPath() {
	siter := NewSiter(gsys.GetPgSites())
	siter.InitNodeSysPath()
}

// 加载应用服务
func (i *InitService) loadAppService() {

}
