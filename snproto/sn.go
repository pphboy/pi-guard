package snproto

type RpcServer interface {
	MonitorServiceServer
	NodeAppServiceServer
	NodeServiceServer
}
