package admin

// login 用jwt
// create
// list
// delete

import (
	"go-ctrl/models"
	rest "pi-rest"

	"github.com/gin-gonic/gin"
)

func NewHttp(group *gin.RouterGroup) {
	s := managerHttp{
		manager: NewManager(),
	}
	group.GET("list", s.list)
	group.POST("delete", s.delete)
	group.POST("", s.create)
	group.POST("login", s.login)
	group.POST("update", s.update)
}

type managerHttp struct {
	manager *Manager
}

func (s *managerHttp) update(c *gin.Context) {
	p := &models.PiManager{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	if err := s.manager.Update(p); err != nil {
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

func (s *managerHttp) login(c *gin.Context) {
	p := &models.PiManager{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	token, err := s.manager.Login(p)
	if err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "login result",
		Data: token,
	})

}

func (s *managerHttp) create(c *gin.Context) {
	p := &models.PiManager{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	if err := s.manager.Create(p); err != nil {
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

func (s *managerHttp) delete(c *gin.Context) {
	p := &models.PiManager{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	if err := s.manager.Delete(p.ManagerId); err != nil {
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

func (s *managerHttp) list(c *gin.Context) {
	l, err := s.manager.List()
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
