package centers

import (
	"go-ctrl/models"
	rest "pi-rest"

	"github.com/gin-gonic/gin"
)

func NewProjectHttp(g *gin.RouterGroup, pm *ProjectManager) *ProjectHttp {
	p := &ProjectHttp{
		manager: pm,
	}
	p.ServeHttp(g)
	return p
}

type ProjectHttp struct {
	manager *ProjectManager
}

func (p *ProjectHttp) ServeHttp(g *gin.RouterGroup) {
	g.GET("list", p.list)
	g.DELETE("", p.delete)
	g.POST("", p.create)
}

func (p *ProjectHttp) list(ctx *gin.Context) {
	list, err := p.manager.List()
	if err != nil {
		// ctx.JSON(500,)
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, &rest.SourceResult{
		Data: list,
		Code: 0,
	})
}

func (p *ProjectHttp) create(ctx *gin.Context) {
	// json bind
	param := &models.PiProject{}

	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	if err := p.manager.Create(param); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "创建项目成功",
	})
}

func (p *ProjectHttp) delete(ctx *gin.Context) {
	param := &models.PiProject{}

	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	if err := p.manager.Delete(param); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "项目删除成功",
	})
}
