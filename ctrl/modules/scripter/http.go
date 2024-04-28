package scripter

import (
	"go-ctrl/models"
	rest "pi-rest"

	"github.com/gin-gonic/gin"
)

func NewScripterHttp(group *gin.RouterGroup) {
	s := scripterHttp{
		sm: NewScripterManager(),
	}
	group.GET("list", s.list)
	group.POST("delete", s.delete)
	group.POST("", s.create)
}

type scripterHttp struct {
	sm *ScripterManager
}

func (s *scripterHttp) create(c *gin.Context) {
	p := &models.PiScript{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	if err := s.sm.Create(p); err != nil {
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

func (s *scripterHttp) delete(c *gin.Context) {
	p := &models.PiScript{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	if err := s.sm.Delete(p.ScriptId); err != nil {
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

func (s *scripterHttp) list(c *gin.Context) {
	l, err := s.sm.List()
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
