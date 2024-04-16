package main

import (
	"flag"
	"go-ctrl/http"

	"github.com/sirupsen/logrus"
)

var (
	dataPath = flag.String("path", "./", "store data path")
	port     = flag.Int("port", 7431, "program running port")
)

func main() {
	logrus.Println("start running")
	s := http.NewCtrlHttp(*port)
	// s.RouterGroup()

	if err := s.Run(); err != nil {
		logrus.Fatal("ctrl http,", err)
	}
}
