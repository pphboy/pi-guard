package centers

import (
	"context"
	"encoding/json"
	"go-ctrl/db"
	cm "go-ctrl/models"
	"go-node/models"
	"go-node/tool"
	rest "pi-rest"
	"snproto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewNodeClientHttp(client *NodeGrpcClient, nodeInfo *cm.PiNode, group *gin.RouterGroup) *NodeClientHttp {
	nch := &NodeClientHttp{
		client:   client,
		nodeInfo: nodeInfo,
	}
	nch.ServeHttp(group.Group("node"))
	return nch
}

// 作用
//   - 包装一个NodeGrpc然后将Grpc的接口以 http的形式暴露
// 类似于 /node/ID/infos
// 这个模块的作用就是，注册为一个Node结点的服务

type NodeClientHttp struct {
	// GetInfoPacket(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MonitorPacket, error)
	// GetNodeInfo(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeSys, error)

	// // 关闭，理论上来说，关闭之后就再也打不开了，
	// Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	// // reboot的作用应该是重启服务器j
	// Reboot(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)

	// 关于App先不写多，到展示的那一步，再进行写多一点
	// LoadAppList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeAppArray, error)
	// UninstallApp(ctx context.Context, in *NodeAppInfo, opts ...grpc.CallOption) (*Result, error)
	// InstallApp(ctx context.Context, in *PiCloudApp, opts ...grpc.CallOption) (*Result, error)
	// Stop(ctx context.Context, in *NodeAppInfo, opts ...grpc.CallOption) (*Result, error)
	// Restart(ctx context.Context, in *NodeAppInfo, opts ...grpc.CallOption) (*Result, error)
	client   *NodeGrpcClient
	nodeInfo *cm.PiNode
}

// 注册路由
func (n *NodeClientHttp) ServeHttp(g *gin.RouterGroup) {
	b := g.Group(n.nodeInfo.NodeId)
	logrus.Infof("node %s register : %s\n", n.nodeInfo.NodeId, b.BasePath())
	b.GET("", n.getNodeInfo)
	b.GET("/monitor", n.monitor)
	b.GET("/applist", n.appList)
	b.GET("/reboot", n.reboot)
	b.GET("/log/:num", n.log)
	b.POST("/install", n.installApp)
	b.POST("/stop", n.stopApp)
	b.POST("/start", n.startApp)
	b.POST("/uninstall", n.uninstallApp)
}

func (n *NodeClientHttp) getNodeInfo(ctx *gin.Context) {
	// return n.nodeInfo
	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "获取成功",
		Data: n.nodeInfo,
	})
}

func (n *NodeClientHttp) monitor(ctx *gin.Context) {
	tmctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	mp, err := n.client.GetInfoPacket(tmctx, &snproto.Empty{})
	if err != nil {
		ctx.JSON(500, rest.SourceResult{
			Msg:  err.Error(),
			Code: 500,
		})
		return
	}
	mpj := &models.MonitorPacket{}

	if err := json.Unmarshal([]byte(mp.Json), mpj); err != nil {
		ctx.JSON(500, rest.SourceResult{
			Msg:  err.Error(),
			Code: 500,
		})
		return
	}
	// return n.nodeInfo
	ctx.JSON(200, rest.SourceResult{
		Data: mpj,
		Msg:  "成功",
		Code: 0,
	})
}

func (n *NodeClientHttp) appList(ctx *gin.Context) {
	tmctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	alist, err := n.client.LoadAppList(tmctx, &snproto.Empty{})
	if err != nil {
		ctx.JSON(500, rest.SourceResult{
			Msg:  err.Error(),
			Code: 500,
		})
		return
	}
	a := []*models.NodeApp{}
	for _, v := range alist.Values {
		a = append(a, tool.ConvertMsgToNodeInfo(v))
	}

	ctx.JSON(200, rest.SourceResult{
		Data: a,
		Msg:  "成功",
		Code: 0,
	})
}

func (n *NodeClientHttp) log(ctx *gin.Context) {
	ru := &rest.RequestUriNum{}
	if err := ctx.ShouldBindUri(ru); err != nil {
		ctx.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	resp, err := n.client.Log(GetOneMinuteCtx(), &snproto.LogParam{
		LineNum: int64(ru.Num),
	})
	if err != nil {
		ctx.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "重启成功，正在重启中",
		Data: resp.Log,
	})
}

func (n *NodeClientHttp) reboot(ctx *gin.Context) {
	c, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	_, err := n.client.Reboot(c, &snproto.Empty{})
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "重启成功，正在重启中",
	})
}

func (n *NodeClientHttp) installApp(ctx *gin.Context) {
	pca := &models.PiCloudApp{}
	if err := ctx.ShouldBindJSON(pca); err != nil {
		ctx.JSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	// pca.AppSite
	res, err := n.client.InstallApp(c, pca.GrpcMsg())
	if err != nil {
		ctx.AbortWithStatusJSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
		Msg:  "安装成功",
		Data: res,
	})
}

type NodeManager struct {
	db *gorm.DB
}

func NewNodeManager() *NodeManager {

	return &NodeManager{
		db: db.DB(),
	}
}

func (nm *NodeManager) List(pid *int) (pns []*cm.PiNode, err error) {
	if err := nm.db.Where("project_id = ?", pid).Find(&pns).Error; err != nil {
		return nil, err
	}
	return
}

func (nm *NodeManager) Get(id string) (*cm.PiNode, error) {
	pn := &cm.PiNode{}
	if err := nm.db.Where("node_id = ?", id).First(pn).Error; err != nil {
		return nil, err
	}

	return pn, nil
}

func (nm *NodeManager) FirstOrCreate(node *cm.PiNode) error {
	if err := nm.db.FirstOrCreate(node).Error; err != nil {
		return err
	}
	return nil
}

func (n *NodeClientHttp) uninstallApp(ctx *gin.Context) {
	nodeApp := models.NodeApp{}

	if err := ctx.ShouldBindJSON(&nodeApp); err != nil {
		ctx.AbortWithStatusJSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	if _, err := n.client.UninstallApp(GetOneMinuteCtx(), nodeApp.Message()); err != nil {
		ctx.AbortWithStatusJSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
	})
}

func (n *NodeClientHttp) stopApp(ctx *gin.Context) {

	nodeApp := models.NodeApp{}
	if err := ctx.ShouldBindJSON(&nodeApp); err != nil {
		ctx.AbortWithStatusJSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	logrus.Infof("start: %+v", nodeApp)

	if _, err := n.client.Stop(GetOneMinuteCtx(), nodeApp.Message()); err != nil {
		ctx.AbortWithStatusJSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(200, &rest.SourceResult{
		Code: 0,
	})
}

func (n *NodeClientHttp) startApp(ctx *gin.Context) {
	nodeApp := models.NodeApp{}

	if err := ctx.ShouldBindJSON(&nodeApp); err != nil {
		ctx.AbortWithStatusJSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	logrus.Infof("start: %+v", nodeApp)

	res, err := n.client.Start(GetOneMinuteCtx(), nodeApp.Message())
	if err != nil {
		ctx.AbortWithStatusJSON(500, &rest.SourceResult{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(200, res)

}

func GetOneMinuteCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	return ctx
}
