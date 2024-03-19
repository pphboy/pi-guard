package service

import (
	"os"
	"syscall"
	"testing"
)

func TestDirectoryExists(t *testing.T) {
	_, err := os.Stat("/tmasdfp")
	t.Log(err)
	t.Log(os.IsNotExist(err))
}

func TestMkdir(t *testing.T) {
	t.Log("start")
	defer t.Log("over")

	m := syscall.Umask(0002)
	defer syscall.Umask(m)

	if err := os.MkdirAll("./pg/db", 0755); err != nil {
		t.Log("pg db", err)
	}

}
