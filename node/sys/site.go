package sys

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
			Name: "Node Database Path",
		},
		{
			Path: root + "/app",
			Name: "Node Database Path",
		},
		{
			Path: root + "/app",
			Name: "Node Database Path",
		},
	}
}

func GetPgSites() []PgSysSite {
	return pgSites
}
