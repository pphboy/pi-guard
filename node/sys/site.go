package sys

import (
	"fmt"
	"os/exec"
)

type PgSysSite struct {
	Path string
	Name string
}

var pgSites []PgSysSite

func init() {
	root := "./pg"
	// 初始化，目录这些
	pgSites = []PgSysSite{
		{
			Path: root,
			Name: "Node Root Path",
		},
		{
			Path: root + "/db",
			Name: "Node Database Path",
		},
		{
			Path: root + "/plugins",
			Name: "Node Plugins Path",
		},
		{
			Path: root + "/app",
			Name: "Node App Path",
		},
		{
			Path: root + "/logs",
			Name: "Node Logs Path",
		},
		{
			Path: root + "/.trash",
			Name: "Node Trash Path",
		},
	}
}

func GetPgSites() []PgSysSite {
	return pgSites
}

type TypePg = int

const (
	PG_ROOT    TypePg = 0
	PG_DB      TypePg = 1
	PG_PLUGINS TypePg = 2
	PG_APP     TypePg = 3
	PG_LOGS    TypePg = 4
	PG_TRASH   TypePg = 5
)

func PgSite(tp TypePg) PgSysSite {
	return pgSites[tp]
}

func KillProcessByPort(port string) error {
	cmd := fmt.Sprintf("lsof -i:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", port)
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		return err
	}
	return nil
}
