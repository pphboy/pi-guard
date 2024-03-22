package tests

import (
	"context"
	"go-node/models"
	"go-node/modules"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"testing"
	"time"
)

func gg() (string, error) {

	resp, err := http.Get("http://127.0.0.1:8081")
	if err != nil {
		log.Println("http get", err)
		return "", err
	}
	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("io read all", err)
		return "", err
	}
	log.Printf("body: %s", bd)
	defer resp.Body.Close()

	return string(bd), nil
}

func TestRunnerApp(t *testing.T) {

	ra := modules.NewRunnerApp(exec.Command("./pghttp"),
		&models.NodeApp{})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ra.Start()
		wg.Done()
	}()

	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second) // 请求客户端
		res, err := gg()
		if err != nil {
			log.Fatal("http get, err:", err)
		}

		t.Log("http get:", res)
		cancel()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		// 10 秒之后关闭服务器
		err := ra.Close()
		if err != nil {
			log.Println("stop err:", err)
		}
		t.Log("close pghttp")
		wg.Done()
	}()

	wg.Add(1)
	c1, cc2 := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		time.Sleep(10 * time.Second)

		res, err := gg()
		t.Log("resource:", res)
		t.Log("http response err,", err)
		if err == nil {
			log.Fatal("err 必须不为nil，访问必须出错")
		}
		cc2()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-c1.Done()

		t.Log("restart")
		ra.Restart()
		time.Sleep(1 * time.Second)
		res, err := gg()
		if err != nil {
			log.Fatal("http get, err:", err)
		}

		t.Log("http get:", res)
		wg.Done()
	}()

	wg.Wait()
}

func TestServiceManager(t *testing.T) {
	t.Log("黄帝崩，藏桥山")
	mg := modules.NewAppDirector()

	var wg sync.WaitGroup

	na := &models.NodeApp{
		NodeAppId:   "ggasdf12",
		NodeAppName: "pghttp",
	}

	ee := exec.Command("./pghttp")

	appSer := modules.NewRunnerApp(ee, na)
	mg.AddService(appSer)

	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Second)

		res, err := gg()
		if err != nil {
			log.Fatal("http get, err:", err)
		}

		t.Log("http get:", res)
		wg.Done()

		if err := appSer.Close(); err != nil {
			log.Fatal("appSer close ", err)
		}

		cancel()
	}()

	wg.Add(1)
	go func() {
		select {
		case <-ctx.Done():
		}
		time.Sleep(3 * time.Second)

		res, err := gg()
		if err == nil {
			log.Fatal("http get, err:", err)
		}

		log.Println("closed err", err)

		t.Log("http get:", res)
		wg.Done()
	}()

	wg.Wait()
}
