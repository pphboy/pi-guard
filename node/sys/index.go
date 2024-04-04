package sys

type Node interface {
	Init()
	initService()
	initGrpcServer()
}

type NodeImpl struct {
}

func NewNode() Node {
	return &NodeImpl{}
}

func (n *NodeImpl) Init() {

}

// service manager
func (n *NodeImpl) initService() {

}

func (n *NodeImpl) initGrpcServer() {

}
