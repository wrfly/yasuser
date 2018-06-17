package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/utils"
	yaml "gopkg.in/yaml.v2"
)

type SrvConfig struct {
	Prefix string `default:"https://u.kfd.me"`
	Port   int    `default:"8084"`
}

type StoreConfig struct {
	DBPath string `default:"./yasuser.db"`
	DBType string `default:"bolt"`
	Redis  string `default:"redis://localhost:6379"`
}

type ShortenerConfig struct {
	Store StoreConfig
}

type Config struct {
	Debug     bool `default:"false"`
	Shortener ShortenerConfig
	Server    SrvConfig
}

func New() *Config {
	conf := Config{
		Server: SrvConfig{},
		Shortener: ShortenerConfig{
			Store: StoreConfig{},
		},
	}
	return &conf
}

func (c *Config) Parse(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		logrus.Fatal(utils.AddLineNum(err))
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.Fatal(utils.AddLineNum(err))
	}

	err = yaml.Unmarshal(bs, c)
	if err != nil {
		logrus.Fatal(utils.AddLineNum(err))
	}
}

func (c *Config) CombineWithENV() {
	utils.ParseConfigEnv(c, []string{"YASUSER"})
}

func (c *Config) EnvConfigLists() []string {
	return utils.EnvConfigLists(c, []string{"YASUSER"})
}

func (c *Config) setDefault() {
	utils.DefaultConfig(c)
}

func (c *Config) Example() {
	c.setDefault()
	bs, err := yaml.Marshal(*c)
	if err != nil {
		logrus.Fatalf("marshal yaml error: %s", utils.AddLineNum(err))
	}
	fmt.Printf("%s", bs)
}
