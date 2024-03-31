package sys

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

const (
	ROOT_DOMAIN = "pi.g"
)

var BASE_DIR = ""

type PgSysSite struct {
	Path string
	Name string
}

func GetPgSites(base string) []PgSysSite {
	BASE_DIR = base
	root := filepath.Join(base, "pg")
	pgSites := []PgSysSite{
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
	pgSites := GetPgSites(BASE_DIR)
	return pgSites[tp]
}

func KillProcessByPort(port string) error {
	cmd := fmt.Sprintf("lsof -i:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", port)
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		return err
	}
	return nil
}
