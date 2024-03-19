package service

import (
	"errors"
	"go-node/models"
	"go-node/service"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
)

func startMockFileServer() {
	http.ListenAndServe(":8081", http.FileServer(http.Dir("./")))
}

func TestCreateTmpFile(t *testing.T) {
	f, err := os.CreateTemp("/tmp", "pgdown")
	t.Log(f.Name())
	if err != nil {
		t.Log(err)
		return
	}
	f.WriteString("原神")

	s, _ := os.Open("/tmp")

	n, _ := s.Readdirnames(100)
	t.Log(n)
	f.Close()

	// os.Remove(f.Name())
}

func TestRmA(t *testing.T) {
	os.RemoveAll("./gg")
}

func TestFileS(t *testing.T) {
	var wg sync.WaitGroup
	go startMockFileServer()
	wg.Add(1)

	go func() {
		resp, err := http.Get("http://127.1:8081/PgHttp.pkg")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Content-Length", resp.ContentLength)
		log.Printf("%+v", resp)
		wg.Done()
	}()

	wg.Wait()
}

func TestPkgInstall(t *testing.T) {
	initer := service.NewIniter()

	p := service.NewPkgService(service.BaseService{
		DB: initer.GetDB(),
	})

	go startMockFileServer()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {

		if err := p.InstallApp(&models.PiCloudApp{
			AppName:    "pghttp",
			AppVersion: 2,
			AppId:      "UUID_123",
			AppSite:    "http://127.1:8081/PgHttp_v2.pkg",
		}); err != nil {
			log.Println("err,", err)
		}

		if err := p.InstallApp(&models.PiCloudApp{
			AppName:    "pghttp",
			AppVersion: 3,
			AppId:      "UUID_123",
			AppSite:    "http://127.1:8081/PgHttp_v3.pkg",
		}); err != nil {
			log.Println("err,", err)
		}

		wg.Done()
	}()

	wg.Wait()
}

func TestUninstallApp(t *testing.T) {
	initer := service.NewIniter()
	p := service.NewPkgService(service.BaseService{
		DB: initer.GetDB(),
	})

	if err := p.UninstallApp(&models.NodeApp{
		NodeAppId:   "412D49368D6C41368A5ADAA9D377BB68",
		NodeAppName: "pghttp",
	}); err != nil {
		t.Fatal("uninstall app", err)
	}

}

func TestLoadAppList(t *testing.T) {

	initer := service.NewIniter()
	p := service.NewPkgService(service.BaseService{
		DB: initer.GetDB(),
	})
	n, err := p.LoadAppList()
	if err != nil {
		if !errors.Is(err, service.ErrLoadAppNotExist) {
			t.Fatal(err)
		}
	}
	log.Printf("%+v", n[0])
}
