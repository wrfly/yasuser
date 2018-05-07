package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"

	"github.com/wrfly/short-url/config"
	"github.com/wrfly/short-url/handler"
	"github.com/wrfly/short-url/handler/db"
	"github.com/wrfly/short-url/routes"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	conf := config.Config{}
	appFlags := []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Usage:       "port number",
			EnvVars:     []string{"PORT"},
			Value:       8080,
			Destination: &conf.Port,
		},
		&cli.StringFlag{
			Name:        "prefix domain",
			Usage:       "short URL prefix",
			EnvVars:     []string{"PREFIX"},
			Value:       "https://u.kfd.me",
			Destination: &conf.Prefix,
		},
		&cli.StringFlag{
			Name:        "db-path",
			Usage:       "database path",
			EnvVars:     []string{"DB_PATH"},
			Value:       "short-url.db",
			Destination: &conf.DBPath,
		},
		&cli.StringFlag{
			Name:        "db-type",
			Usage:       "database type: redis or file",
			EnvVars:     []string{"DB_TYPE"},
			Value:       "file",
			Destination: &conf.DBType,
		},
		&cli.StringFlag{
			Name:        "redis",
			Usage:       "database path",
			EnvVars:     []string{"REDIS"},
			Value:       "localhost:6379/0",
			Destination: &conf.Redis,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "log level: debug",
			EnvVars:     []string{"DEBUG"},
			Value:       false,
			Destination: &conf.Debug,
		},
	}

	app := &cli.App{
		Name:    "short-url",
		Usage:   "short your url",
		Authors: author,
		Version: fmt.Sprintf("Version: %s\tCommit: %s\tDate: %s",
			Version, CommitID, BuildAt),
		Flags: appFlags,
		Action: func(c *cli.Context) error {
			var (
				shorterDB db.Database
				err       error
			)
			switch conf.DBType {
			case "file":
				shorterDB, err = db.NewBoltDB(conf.DBPath)
				if err != nil {
					logrus.Fatal(err)
				}
			case "redis":
				redisHosts := strings.Split(conf.Redis, ",")
				logrus.Info(redisHosts)
				return nil
				// shorterDB, err = db.NewDB(conf.DBPath)
				// if err != nil {
				// 	logrus.Fatal(err)
				// }
			default:
				logrus.Fatalf("unknown db type: %s", conf.DBType)
			}

			shorter := handler.Shorter{
				DB: shorterDB,
			}

			if conf.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				gin.SetMode(gin.ReleaseMode)
			}

			return routes.Serve(&conf, &shorter)
		},
	}

	app.Run(os.Args)
}
