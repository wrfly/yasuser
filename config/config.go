package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/ecp"
	yaml "gopkg.in/yaml.v2"

	"github.com/wrfly/yasuser/utils"
)

type SrvConfig struct {
	Domain string `default:"https://u.kfd.me"`
	Port   int    `default:"8084"`
	Limit  int64  `default:"10"`
	Pprof  bool   `default:"false"`
	GAID   string `default:"62244864-8"`
}

type StoreConfig struct {
	DBPath string `default:"./yasuser.db"`
	DBType string `default:"bolt"`
	Redis  string `default:"redis://localhost:6379"`
}

type list struct {
	WhiteList []string
	BlackList []string
}

type Filter struct {
	Domain  list
	Keyword list
}

type Config struct {
	Debug  bool   `default:"false"`
	Auth   string `default:"password"`
	Store  StoreConfig
	Server SrvConfig
	Filter Filter
}

func New() *Config {
	return &Config{}
}

func (c *Config) Parse(filePath string) {
	if filePath == "" {
		return
	}
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

func (c *Config) Example() {
	ecp.Parse(c)

	bs, err := yaml.Marshal(*c)
	if err != nil {
		logrus.Fatalf("marshal yaml error: %s", utils.AddLineNum(err))
	}
	fmt.Printf("%s", bs)
}
