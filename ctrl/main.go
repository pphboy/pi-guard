package main

import (
	"flag"
	"go-ctrl/db"
	"go-ctrl/http"
	"go-ctrl/modules/appm"
	centers "go-ctrl/modules/center"
	"os"
	"path"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	dataPath = flag.String("path", "./", "store data path")
	ctrlName = flag.String("name", "ctrl", "controller name")
	port     = flag.Int("port", 7431, "program running port")
	picPath  = path.Join(*dataPath, "static")
)

func main() {
	flag.Parse()
	initOperation()

	logrus.Println("start running")
	s := http.NewCtrlHttp(*port)

	s.RouterGroup("").Static("/assets", picPath)

	// center 中心
	cm := centers.NewManager(s.RouterGroup("project"))

	// Project中心
	pm := centers.NewProjectManager(cm)

	// appStore
	appm.NewAppHttp(picPath, s.RouterGroup("appm"))

	centers.NewProjectHttp(s.RouterGroup("project"), pm)

	s.RouterGroup("project")

	if err := s.Run(); err != nil {
		logrus.Fatal("ctrl http,", err)
	}
}

func initOperation() {
	// static file dir
	_, err := os.Stat(picPath)
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(picPath, 0744); err != nil {
			logrus.Error("mkdir picPath", picPath, err)
		}
	}

	logrus.Print("running main init")
	db.Init(*dataPath, *ctrlName)
}
