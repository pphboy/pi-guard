package service

import (
	"fmt"
	"go-node/sys"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Siter interface {
	InitNodeSysPath()
}

func NewSiter(ps []sys.PgSysSite) Siter {
	return &siter{
		ps: ps,
	}
}

type siter struct {
	ps []sys.PgSysSite
	ne []sys.PgSysSite
}

func (s *siter) InitNodeSysPath() {
	// 获取不存在的路径
	s.getNotExistsPgSite()
	// 创建路径
	s.createPgSites()
}

func (s *siter) getNotExistsPgSite() {
	var res []sys.PgSysSite
	for _, v := range s.ps {
		if !IsFileExist(v.Path) {
			logrus.Warnf("%+v dont existing", v)
			res = append(res, v)
		}
	}
	s.ne = res
}

func (s *siter) createPgSites() error {
	m := syscall.Umask(0002)
	defer syscall.Umask(m)

	// has create error PgSysSite
	errStr := ""
	for _, v := range s.ne {
		if err := os.MkdirAll(v.Path, 0755); err != nil {
			errStr = fmt.Sprintf("%s failed to created: %+v,%v\n", errStr, v, err)
		}
	}

	if len(errStr) > 0 {
		return fmt.Errorf("has failed created list \n%s", errStr)
	} else {
		return nil
	}
}

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDirExist(path string) bool {
	return IsFileExist(path)
}
