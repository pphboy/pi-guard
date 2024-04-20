package main

import (
	"flag"
	"go-ctrl/db"
	"go-ctrl/http"
	centers "go-ctrl/modules/center"

	"github.com/sirupsen/logrus"
)

var (
	dataPath = flag.String("path", "./", "store data path")
	ctrlName = flag.String("name", "ctrl", "controller name")
	port     = flag.Int("port", 7431, "program running port")
)

func init() {
	logrus.Print("running main init")
	db.Init(*dataPath, *ctrlName)
}

func main() {
	logrus.Println("start running")
	s := http.NewCtrlHttp(*port)
	// s.RouterGroup()
	// db := db.DB()

	// center 中心
	cm := centers.NewManager(s.RouterGroup("project"))

	// Project中心
	pm := centers.NewProjectManager(cm)

	centers.NewProjectHttp(s.RouterGroup("project"), pm)

	s.RouterGroup("project")

	if err := s.Run(); err != nil {
		logrus.Fatal("ctrl http,", err)
	}
}
