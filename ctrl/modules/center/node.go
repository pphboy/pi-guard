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
	b.GET("/", n.getNodeInfo)    // 返回当前结点的基本消息
	b.GET("/monitor", n.monitor) // 返回当前结点的基本消息
	b.GET("/applist", n.appList) // 返回当前结点的基本消息
}

func (n *NodeClientHttp) getNodeInfo(ctx *gin.Context) {
	// return n.nodeInfo
	ctx.JSON(200, n.nodeInfo)
}

func (n *NodeClientHttp) monitor(ctx *gin.Context) {
	tmctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	mp, err := n.client.GetInfoPacket(tmctx, &snproto.Empty{})
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
			"code":  500,
		})
		return
	}
	mpj := &models.MonitorPacket{}

	if err := json.Unmarshal([]byte(mp.Json), mpj); err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
			"code":  500,
		})
		return
	}
	// return n.nodeInfo
	ctx.JSON(200, mpj)
}

func (n *NodeClientHttp) appList(ctx *gin.Context) {
	tmctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	alist, err := n.client.LoadAppList(tmctx, &snproto.Empty{})
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err,
			"code":  500,
		})
		return
	}
	var a []*models.NodeApp
	for _, v := range alist.Values {
		a = append(a, tool.ConvertMsgToNodeInfo(v))
	}

	ctx.JSON(200, rest.SourceResult{
		Data: a,
		Msg:  "成功",
		Code: 0,
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