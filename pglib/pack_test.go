package pglib

import (
	"log"
	"os/exec"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPkg(t *testing.T) {
	// viper.SetDefault()

}

func TestViper(t *testing.T) {

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.SetConfigFile("testpkg/pkg.toml")

	err := v.ReadInConfig()

	if err != nil {
		panic(err)
	}
	Conf := Config{}

	if err := v.Unmarshal(&Conf); err != nil {
		log.Fatal(err, "Unmarshal")
	}

	log.Print(v.Get("PkgConfig.Name"))
	log.Printf("%+v", Conf)
}

func TestLoadConfigFile(t *testing.T) {
	t.Run("test load config File", func(t *testing.T) {
		a := assert.New(t)

		cfg, err := LoadPackageConfig("./testpkg/pkg.toml")
		if err != nil {
			t.Fatal(err)
		}

		if a.NotNil(cfg, "config never load nil") {
			t.Logf("load config succeed! \n%+v", cfg)
		}
	})
}

func TestCmd(t *testing.T) {
	bs, _ := exec.Command("ls").Output()
	log.Printf("\n%s", bs)
}

func TestPack(t *testing.T) {
	cfg, _ := LoadPackageConfig("./testpkg/pkg.toml")
	PackPkg(cfg, "testpkg")
}

func TestPackPiguardPakcage(t *testing.T) {
	a := assert.New(t)

	err := PackPiguardPackage("./ggg")

	a.Nil(err, "package exception", err)
}

func TestUnpackPkg(t *testing.T) {
	UnpackPkg("./ggg.pkg", "./gg123")
}
