package service

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCom(t *testing.T) {
	ap, err := filepath.Abs("./a.sh")
	if err != nil {
		t.Fatal(err)
	}
	os.Chmod(ap, 0777)

	t.Log(ap)
	cmd := exec.Command("sh", ap)
	go func() {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	p, err := cmd.StdoutPipe()

	if err != nil {
		t.Fatal(err)
	}

	r, err := io.ReadAll(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", r)
}
