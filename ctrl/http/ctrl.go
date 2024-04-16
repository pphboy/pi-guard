package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type CtrlHttp interface {
	// 启动Http服务器
	Run() error
	// 获取指定根目录下的Group目录
	RouterGroup(string) *gin.RouterGroup
}

func NewCtrlHttp(port int) CtrlHttp {
	e := gin.Default()
	p := &piguardCtrlHttp{
		port:   port,
		engine: e,
	}
	return p
}

type piguardCtrlHttp struct {
	engine *gin.Engine
	port   int
	router *gin.RouterGroup
}

func (p *piguardCtrlHttp) Run() error {
	return p.engine.Run(fmt.Sprintf(":%d", p.port))
}

func (p *piguardCtrlHttp) RouterGroup(groupName string) *gin.RouterGroup {
	return p.engine.Group(groupName)
}
