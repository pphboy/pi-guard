package rproxy

import (
	"go-node/service"
	"sync"

	"gorm.io/gorm"
)

// 作用，将所有访问 app*.node.pi.g的域名导向对应的服务
// 比如 pghttp.node.pi.g  -> 127.0.0.1:8080
// 原理是： node.pi.g,appx,node.pi.g -> 192.168.1.xxx
// 总之本节点的解析都会解析到当前服务器，然后由本服务进行处理

type ReverseProxy interface {
	RefreshHosts() error
	HandleAppHost(host string) error
	HostMap() map[string]string
}

func NewRProxyer(initer service.Initer) ReverseProxy {

	return &RProxyer{
		db: initer.GetDB(),
	}
}

type RProxyer struct {
	db      *gorm.DB
	hostMap sync.Map
}

// 将hostMap重新写一个，进行返回
func (r *RProxyer) HostMap() map[string]string {

	return nil
}

func (r *RProxyer) RefreshHosts() error {

	return nil
}

func (r *RProxyer) HandleAppHost(host string) error {

	return nil
}
