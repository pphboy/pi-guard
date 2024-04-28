package ssher

import (
	"go-ctrl/models"
	rest "pi-rest"

	"github.com/gin-gonic/gin"
)

func NewSsherHttp(group *gin.RouterGroup) {
	s := ssherHttp{
		sshM: NewSsherManager(),
	}
	group.GET("list", s.list)
	group.POST("delete", s.delete)
	group.POST("", s.create)
	group.GET("connect", s.connect)
}

type ssherHttp struct {
	sshM *SsherManager
}

func (s *ssherHttp) connect(c *gin.Context) {

}
func (s *ssherHttp) create(c *gin.Context) {
	p := &models.PiTerminal{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	if err := s.sshM.Create(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "create success!",
	})
}

func (s *ssherHttp) delete(c *gin.Context) {
	p := &models.PiTerminal{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	if err := s.sshM.Delete(p.TerminalId); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "delete success!",
	})
}

func (s *ssherHttp) list(c *gin.Context) {
	l, err := s.sshM.List()
	if err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "get list success!",
		Data: l,
	})
}
