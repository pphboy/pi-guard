package pglib

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/spf13/viper"
)

type PkgType string

const (
	PKGFILE_NAME = "pkg.toml"

	TYPE_SYS  = "SYS"
	TYPE_NORM = "NORM"
)

type PkgConfig struct {
	Name  string
	Intro string
	Type  PkgType
	Hash  string
	Port  string
	// 应用下可执行文件名即可，不要使用 ./run.sh , 正确是 run.sh
	Exec    string
	Version *int
}

type Config struct {
	PkgConfig
}

func LoadPackageConfig(filepath string) (cfg *Config, err error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.SetConfigFile(filepath)

	if err != v.ReadInConfig() {
		return nil, err
	}

	if err != v.Unmarshal(&cfg) {
		return nil, err
	}

	return
}

func PackPkg(cfg *Config, dir string) (err error) {
	if cfg.Version == nil {
		return fmt.Errorf("option version is not set")
	}
	execFile := path.Join(dir, cfg.Exec)
	log.Println("Exec:", execFile)
	if !IsFileExist(execFile) {
		return fmt.Errorf("file don't have exec file,%w", os.ErrNotExist)
	}

	cmd := exec.Command("tar", "-zcvf",
		fmt.Sprintf("%s.pkg", cfg.Name), "-C", dir, ".")

	log.Println("cmd:", cmd)

	d, err := cmd.Output()
	if err != nil {
		log.Fatal("cmd output ", err)
	}
	log.Printf("package dir & files :\n%s", d)

	if err != nil {
		return err
	}
	return
}

func PackPiguardPackage(dir string) (err error) {
	// pkg.config exist dir
	cfgFile := path.Join(dir, PKGFILE_NAME)
	// log.Println("config:", cfgFile)
	// validate config of package
	_, err = os.Stat(cfgFile)
	if err == nil {
		// pack a bag
		cfg, err := LoadPackageConfig(cfgFile)
		log.Println("READ Config File Succeed!!", cfg)
		if err != nil || cfg == nil {
			log.Fatal("load package config err,", err, nil)
		}

		if err := PackPkg(cfg, dir); err != nil {
			log.Fatal("package pkg file err,", err)
		}

		log.Println("package completed!!!")
	}
	return
}

func UnpackPkg(filepath string, installSiteDir string) {
	syscall.Umask(0002)

	// fileMode must be 0755,
	err := os.MkdirAll(installSiteDir, 0755)

	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("mkdir", err)
		}
	}
	cmd := exec.Command("tar",
		"-zxvf", filepath, "-C", installSiteDir)

	log.Println("cmd:", cmd)
	log.Println("installDir:", installSiteDir)
	d, err := cmd.Output()

	if err != nil {
		log.Println("failed install , remove dir ")

		if err := os.Remove(installSiteDir); err != nil {
			log.Println("remove file ", err)
		}
		log.Fatal("unpack error ", err)
	}

	log.Printf("unpack logs:\n%s", d)

}
