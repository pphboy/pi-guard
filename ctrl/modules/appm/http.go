package appm

import (
	"fmt"
	"go-node/models"
	"path"
	"pglib"
	rest "pi-rest"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// get: list 获取列表
// post: file  上传文件
// post: create 创建
// post: update 更新

type appMHttp struct {
	store      *AppStore
	picRoot    string
	ctrlDomain string
}

func NewAppHttp(picRoot, ctrlDomain string, g *gin.RouterGroup) *appMHttp {
	am := &appMHttp{
		store:      NewAppStore(),
		picRoot:    picRoot,
		ctrlDomain: ctrlDomain,
	}
	am.serverHttp(g)
	return am
}

func (a *appMHttp) serverHttp(g *gin.RouterGroup) {
	g.GET("/:id", a.findOne)
	g.GET("/list", a.list)
	// 删除app
	g.DELETE("/:id", a.remove)
	// 删除或更新
	g.POST("", a.create)
	// 上传文件
	g.POST("/file", a.upload)
}

// create & update 共用这个接口
func (a *appMHttp) create(c *gin.Context) {
	pc := &models.PiCloudApp{}
	if err := c.ShouldBindJSON(pc); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	var err error = nil
	if strings.Compare(pc.AppId, "") == 0 {
		err = a.store.Create(pc)
	} else {
		err = a.store.Update(pc)
	}

	if err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "modify success!",
	})

}

func (a *appMHttp) findOne(c *gin.Context) {
	rq := &rest.RequestUriId{}
	if err := c.ShouldBindUri(rq); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	p, err := a.store.Get(rq.ID)
	if err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "get success!",
		Data: p,
	})
}

func (a *appMHttp) remove(c *gin.Context) {
	rq := &rest.RequestUriId{}
	if err := c.ShouldBindUri(rq); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	if err := a.store.Delete(rq.ID); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "remove success!",
	})
}

func (a *appMHttp) list(c *gin.Context) {
	p, err := a.store.List()
	if err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "get list",
		Data: p,
	})
}

func (a *appMHttp) upload(c *gin.Context) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	file, _ := c.FormFile("file")

	logrus.Println(file.Filename)
	fn := fmt.Sprintf("%s.pkg", uuid.NewString())
	// 直接是在服务器运行的位置
	dstName := path.Join(a.picRoot, fn)

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, dstName); err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	// 保存文件之后将文件解压到/tmp下，然后拿到版本
	unpackPath := path.Join("/tmp", uuid.NewString())
	pglib.UnpackPkg(dstName, unpackPath)
	cfg, err := pglib.LoadPackageConfig(path.Join(unpackPath, pglib.PKGFILE_NAME))
	if err != nil {
		c.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "upload file",
		Data: gin.H{
			"path":    fmt.Sprintf("http://%s/assets/%s", a.ctrlDomain, fn),
			"version": cfg.Version,
			"name":    cfg.Name,
		},
	})
}
