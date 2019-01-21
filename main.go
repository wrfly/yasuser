package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/ecp"
	"gopkg.in/urfave/cli.v2"

	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/filter"
	"github.com/wrfly/yasuser/routes"
	"github.com/wrfly/yasuser/shortener"
)

var appName = "yasuser"

func main() {

	app := &cli.App{
		Name:    appName,
		Usage:   "Yet another self-hosted URL shortener.",
		Authors: author,
		Version: fmt.Sprintf("Version: %s\tCommit: %s\tDate: %s",
			Version, CommitID, BuildAt),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "config file path",
				Aliases: []string{"c"},
				Value:   "./config.yml",
			},
			&cli.BoolFlag{
				Name:    "example",
				Usage:   "config file example",
				Aliases: []string{"e"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:  "env-list",
				Usage: "config environment lists",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			conf := config.New()
			if c.Bool("example") {
				conf.Example()
				return nil
			}
			if c.Bool("env-list") {
				for _, e := range ecp.List(conf, appName) {
					fmt.Println(e)
				}
				return nil
			}
			conf.Parse(c.String("config"))
			if err := ecp.Parse(conf, appName); err != nil {
				logrus.Error(err)
				return err
			}

			if conf.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				gin.SetMode(gin.ReleaseMode)
			}

			err := routes.Serve(
				conf.Server,
				shortener.New(conf.Store),
				filter.New(conf.Filter),
			)
			if err != nil {
				logrus.Error(err)
			}

			return nil
		},
	}
	app.CustomAppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}

AUTHOR:
	{{range .Authors}}{{ . }}
	{{end}}
VERSION:
	{{.Version}}

OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}`

	app.Run(os.Args)
}
