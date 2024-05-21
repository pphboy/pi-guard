package rproxy

import (
	"errors"
	"fmt"
	"go-node/service"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
)

// 作用，将所有访问 app*.node.pi.g的域名导向对应的服务
// 比如 pghttp.node.pi.g  -> 127.0.0.1:8080
// 原理是： node.pi.g,appx,node.pi.g -> 192.168.1.xxx
// 总之本节点的解析都会解析到当前服务器，然后由本服务进行处理
var (
	ErrDontMatchHost = errors.New("dont match host in this host")
)

type ReverseProxy interface {
	RefreshCache() error
	HostMap() map[string]string
	ReverseHandle(http.ResponseWriter, *http.Request)
}

func NewRProxyer(bs service.BaseService, bus EventBus.Bus) ReverseProxy {
	rp := &RProxyer{
		//		sysDao: dao.NewSysDao(bs.DB),
		pkgService: service.NewPkgService(bs, bus),
	}
	rp.RefreshCache()

	return rp
}

type RProxyer struct {
	pkgService service.PkgService
	// like {"app.node.pi.g": <*httputil.ReverProxy> }
	proxyMap sync.Map
	// like { "app.node.pi.g":":8080"}
	hostMap sync.Map
}

// 将hostMa重新写一个，进行返回
func (r *RProxyer) HostMap() map[string]string {
	hosts := make(map[string]string)

	r.hostMap.Range(func(key, value any) bool {
		hosts[key.(string)] = value.(string)
		return true
	})

	return hosts
}

func (r *RProxyer) RefreshCache() error {
	alist, err := r.pkgService.LoadAppList()

	if err != nil {
		return err
	}
	r.hostMap = sync.Map{}
	// ip在哪里找
	for _, v := range alist {
		u, err := getTargetUrl(v.NodeAppPort)
		if err != nil {
			continue
		}
		// 将所有的反向代理全部存起来
		r.proxyMap.Store(v.NodeAppDomain, httputil.NewSingleHostReverseProxy(u))
		// ":<port>"
		// 就是域名对应的本机的端口
		r.hostMap.Store(v.NodeAppDomain, fmt.Sprintf(":%s", v.NodeAppPort))
	}

	return nil
}

func (r *RProxyer) ReverseHandle(w http.ResponseWriter, req *http.Request) {
	// 拿到请求的host
	hs := strings.Split(req.Host, ":")
	host := hs[0]
	// port := hs[1]
	// port无用
	logrus.Println("host:", host)
	p, ok := r.proxyMap.Load(host)
	if ok {
		ps, ok := p.(*httputil.ReverseProxy)
		if ok {
			ps.ServeHTTP(w, req)
		} else {
			logrus.Error("get exception type of RProxyer.proxyMap")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("reverse proxy type exception"))
		}
	} else {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(ErrDontMatchHost.Error()))
		logrus.Warningf("%v, host: %v", ErrDontMatchHost, req.Host)
	}

}

// 获取对应的targetUrl
func getTargetUrl(port string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("http://127.0.0.1:%s", port))
}
