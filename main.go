package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"

	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/routes"
	"github.com/wrfly/yasuser/shortener"
	"gopkg.in/urfave/cli.v2"
)

func main() {

	conf := config.Config{
		Server: config.SrvConfig{},
		Shortener: config.ShortenerConfig{
			Store: config.StoreConfig{},
		},
	}

	app := &cli.App{
		Name:    "yasuser",
		Usage:   "Yet another self-hosted URL shortener.",
		Authors: author,
		Version: fmt.Sprintf("Version: %s\tCommit: %s\tDate: %s",
			Version, CommitID, BuildAt),
		Action: func(c *cli.Context) error {
			if conf.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				gin.SetMode(gin.ReleaseMode)
			}

			return routes.Serve(conf.Server, shortener.New(conf.Shortener))
		},
	}
	app.CustomAppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}

OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}
AUTHOR:
	{{range .Authors}}{{ . }}
	{{end}}
VERSION:
	{{.Version}}
`
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Aliases:     []string{"p"},
			Usage:       "port number",
			EnvVars:     []string{"PORT"},
			Value:       8080,
			Destination: &conf.Server.Port,
		},
		&cli.StringFlag{
			Name:        "prefix domain",
			Usage:       "short URL prefix",
			EnvVars:     []string{"PREFIX"},
			Value:       "https://u.kfd.me",
			Destination: &conf.Server.Prefix,
		},
		&cli.StringFlag{
			Name:        "db-path",
			Usage:       "database path",
			EnvVars:     []string{"DB_PATH"},
			Value:       "/data/yasuser.db",
			Destination: &conf.Shortener.Store.DBPath,
		},
		&cli.StringFlag{
			Name:        "db-type",
			Usage:       "database type: redis or bolt",
			EnvVars:     []string{"DB_TYPE"},
			Value:       "bolt",
			Destination: &conf.Shortener.Store.DBType,
		},
		&cli.StringFlag{
			Name:        "redis",
			Usage:       "redis host address",
			EnvVars:     []string{"REDIS"},
			Value:       "localhost:6379/0",
			Destination: &conf.Shortener.Store.Redis,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "log level: debug",
			Destination: &conf.Debug,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}
