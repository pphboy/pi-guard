package tests

import (
	"go-node/modules/rproxy"
	"go-node/service"
	"net/http"
	"sync"
	"testing"
)

func TestSyncMap(t *testing.T) {
	var sm sync.Map

	for i := 0; i < 10; i++ {
		sm.Store(i, i)
	}

	sm.Range(func(key, value any) bool {
		t.Logf("k:%d,v:%d", key, value)
		return true
	})

}

// 我想到的是，其实省掉这种使用类型的测试，并省不了时间
// 写测试花不了多少时间
func TestReverProxy(t *testing.T) {
	initer := service.NewIniter("../../../fs_root")
	rp := rproxy.NewRProxyer(service.BaseService{
		DB: initer.GetDB(),
	})
	t.Log("host map", rp.HostMap())

	http.HandleFunc("/", rp.ReverseHandle)
	// 测试基于 80
	// host : ndev.pi.g 绑定 127.0.0.1
	// 先将pkg装好

	// 用service将app全部启动
	// 启动代理服务器
	// 访问 appx.node.pi.g
	t.Fatal(http.ListenAndServe(":8090", nil))

}
