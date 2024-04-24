package nodegrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"go-node/models"
	"go-node/modules/service"
	gs "go-node/service"
	"go-node/tool"
	"log"
	"net"
	"snproto"
	sp "snproto"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NodeRpcService struct {
	sp.UnimplementedMonitorServiceServer
	sp.UnimplementedNodeAppServiceServer
	sp.UnimplementedNodeServiceServer

	monitor gs.NodeMonitor
	pkg     gs.PkgService
	sys     gs.SysService
	// 这里还需要扩展一个nodeApp serviceManager
	//看是将Node放到这里还是将其他的内容放到这里
	sm service.ServiceManager
}

// serviceManager属于核心模块，不能直接new
func NewNodeRpcService(sm service.ServiceManager, bs gs.BaseService) *NodeRpcService {
	return &NodeRpcService{
		monitor: gs.NewNodeMonitor(1600, 10*time.Second, func(mp []*models.MonitorPacket) {
			logrus.Printf("monitor packet,%+v", mp)
		}),
		pkg: gs.NewPkgService(bs),
		sys: gs.NewSysService(bs),
		sm:  sm,
	}
}

// nodeService
func (n *NodeRpcService) GetNodeInfo(context.Context, *sp.Empty) (*sp.NodeSys, error) {

	nf, err := n.sys.GetSysInfo()
	if err != nil {
		return nil, err
	}

	ns := &sp.NodeSys{
		NodeId:     nf.NodeId,
		NodeName:   nf.NodeName,
		NodeStatus: int32(nf.NodeStatus),
		NodeDomain: nf.NodeDomain,
		CreatedAt:  timestamppb.New(*nf.CreatedAt),
		UpdatedAt:  timestamppb.New(*nf.UpdatedAt),
	}

	if nf.DeletedAt != nil {
		ns.DeletedAt = timestamppb.New(*nf.DeletedAt)
	}

	return ns, nil

}
func (n *NodeRpcService) Shutdown(context.Context, *sp.Empty) (*sp.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}
func (n *NodeRpcService) Reboot(context.Context, *sp.Empty) (*sp.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reboot not implemented")
}

// pkgService
func (n *NodeRpcService) LoadAppList(context.Context, *sp.Empty) (*sp.NodeAppArray, error) {
	ns, err := n.pkg.LoadAppList()
	if err != nil {
		return nil, err
	}

	var arrAny []*snproto.NodeAppInfo
	for _, v := range ns {
		arrAny = append(arrAny, v.Message())
	}

	return &sp.NodeAppArray{
		Values: arrAny,
	}, nil
}
func (n *NodeRpcService) UninstallApp(ctx context.Context, na *sp.NodeAppInfo) (*sp.Result, error) {

	if err := n.pkg.UninstallApp(&models.NodeApp{
		NodeAppId:   na.NodeAppId,
		NodeAppName: na.NodeAppName,
	}); err != nil {
		return nil, err
	}

	ad, err := anypb.New(na)
	if err != nil {
		return nil, err
	}
	return &sp.Result{
		Code:    snproto.DB_DELETE_SUCCEED,
		Message: fmt.Sprintf("删除{%s}应用成功", na.NodeAppName),
		Data:    ad,
	}, nil
}

func (n *NodeRpcService) InstallApp(ctx context.Context, spc *sp.PiCloudApp) (*sp.Result, error) {
	pc := &models.PiCloudApp{}
	if err := n.pkg.InstallApp(pc.ResolveGrpcMsg(spc)); err != nil {
		return nil, err
	}

	sa, err := anypb.New(spc)
	if err != nil {
		return nil, err
	}
	return &sp.Result{
		Code:    snproto.DB_DELETE_SUCCEED,
		Message: fmt.Sprintf("安装[%s]{%s}应用成功", spc.AppId, spc.AppName),
		Data:    sa,
	}, nil
}

func (n *NodeRpcService) Stop(ctx context.Context, nai *sp.NodeAppInfo) (*sp.Result, error) {
	ra, err := n.sm.GetServiceByApp(tool.ConvertRpcInfoToNodeInfo(nai))
	if err != nil {
		return nil, err
	}

	if err := ra.Close(); err != nil {
		return nil, err
	}

	return &sp.Result{
		Code:    snproto.APP_CLOSED_SUCCEED,
		Message: fmt.Sprintf("关闭[%s]{%s}应用成功", nai.NodeAppId, nai.NodeAppName),
		Data:    nil,
	}, nil
}

func (n *NodeRpcService) Restart(ctx context.Context, nai *sp.NodeAppInfo) (*sp.Result, error) {
	ra, err := n.sm.GetServiceByApp(tool.ConvertRpcInfoToNodeInfo(nai))
	if err != nil {
		return nil, err
	}

	if err := ra.Restart(); err != nil {
		return nil, err
	}

	return &sp.Result{
		Code:    snproto.APP_RESTART_SUCCEED,
		Message: fmt.Sprintf("重启[%s]{%s}应用成功", nai.NodeAppId, nai.NodeAppName),
		Data:    nil,
	}, nil
}

func (n *NodeRpcService) Start(ctx context.Context, nai *sp.NodeAppInfo) (*sp.Result, error) {
	ra, err := n.sm.GetServiceByApp(tool.ConvertRpcInfoToNodeInfo(nai))
	if err != nil {
		return nil, err
	}

	if err := ra.Start(); err != nil {
		return nil, err
	}

	return &sp.Result{
		Code:    snproto.APP_RESTART_SUCCEED,
		Message: fmt.Sprintf("启动[%s]{%s}应用成功", nai.NodeAppId, nai.NodeAppName),
		Data:    nil,
	}, nil
}

// monitorService
func (n *NodeRpcService) GetInfoPacket(context.Context, *sp.Empty) (*sp.MonitorPacket, error) {
	p, err := n.monitor.GetInfoPacket()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	mp := &sp.MonitorPacket{
		Json: string(b),
	}

	return mp, nil
}

func RunRpcServer(port int, ns snproto.RpcServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalf("failed to listen %d", port)
	}

	s := grpc.NewServer()

	snproto.RegisterMonitorServiceServer(s, ns)
	snproto.RegisterNodeAppServiceServer(s, ns)
	snproto.RegisterNodeServiceServer(s, ns)

	logrus.Printf("listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
