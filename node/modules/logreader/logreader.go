package logreader

import (
	"bufio"
	"bytes"
	"fmt"
	"go-node/models"
	"go-node/sys"
	"os"
	"path/filepath"
)

type LogReader interface {
	Log(lines int) (*bytes.Buffer, error)
}

func NewLogReader(n *models.NodeSys) LogReader {
	return &nodeLogReader{
		nodeInfo: n,
	}
}

type nodeLogReader struct {
	nodeInfo *models.NodeSys
}

func (n *nodeLogReader) Log(lines int) (*bytes.Buffer, error) {
	niLog := filepath.Join(sys.PgSite(sys.PG_LOGS).Path, fmt.Sprintf("%s_sys.log", n.nodeInfo.NodeName))
	bf := bytes.Buffer{}
	f, err := os.Open(niLog)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := bufio.NewScanner(f)

	for c.Scan() && lines >= 0 {
		bf.WriteString(c.Text() + "\n")
		lines--
	}

	return &bf, nil
}
