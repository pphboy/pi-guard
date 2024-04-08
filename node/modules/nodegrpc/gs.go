package nodegrpc

import (
	"context"
	"fmt"
	"go-node/models"
	"go-node/service"
	"snproto"
	sp "snproto"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NodeRpcService struct {
	sp.UnimplementedMonitorServiceServer
	sp.UnimplementedNodeAppServiceServer
	sp.UnimplementedNodeServiceServer

	monitor service.NodeMonitor
	pkg     service.PkgService
	sys     service.SysService
	// 这里还需要扩展一个nodeApp serviceManager
	// 看是将Node放到这里还是将其他的内容放到这里
}

func NewNodeRpcService(bs service.BaseService) *NodeRpcService {
	return &NodeRpcService{
		monitor: service.NewNodeMonitor(1600, 10*time.Second, func(mp []*models.MonitorPacket) {
			logrus.Printf("monitor packet,%+v", mp)
		}),
		pkg: service.NewPkgService(bs),
		sys: service.NewSysService(bs),
	}
}

// nodeService
func (n *NodeRpcService) GetNodeInfo(context.Context, *sp.Empty) (*sp.NodeSys, error) {

	nf, err := n.sys.GetSysInfo()
	if err != nil {
		return nil, err
	}

	return &sp.NodeSys{
		NodeId:     nf.NodeId,
		NodeName:   nf.NodeName,
		NodeStatus: int32(nf.NodeStatus),
		NodeDomain: nf.NodeDomain,
		CreatedAt:  timestamppb.New(*nf.CreatedAt),
		UpdatedAt:  timestamppb.New(*nf.UpdatedAt),
		DeletedAt:  timestamppb.New(*nf.DeletedAt),
	}, nil
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

	var arrAny []*anypb.Any
	for _, v := range ns {
		d, err := anypb.New(v.Message())
		if err != nil {
			return nil, err
		}
		arrAny = append(arrAny, d)
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
func (n *NodeRpcService) Stop(context.Context, *sp.NodeAppInfo) (*sp.Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (n *NodeRpcService) Restart(context.Context, *sp.NodeAppInfo) (*sp.Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Restart not implemented")
}

// monitorService
func (n *NodeRpcService) GetInfoPacket(context.Context, *sp.Empty) (*sp.MonitorPacket, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfoPacket not implemented")
}

func NewGrpcServer() {

}
