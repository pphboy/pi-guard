package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestComg(t *testing.T) {
	ap, err := filepath.Abs("./pghttp")
	if err != nil {
		t.Fatal(err)
	}
	os.Chmod(ap, 0777)

	t.Log(ap)
	cmd := exec.Command(ap, "&")
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)
}
