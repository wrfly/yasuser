package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/wrfly/short-url/config"
	"github.com/wrfly/short-url/handler"
	"github.com/wrfly/short-url/handler/db"
	"github.com/wrfly/short-url/routes"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	conf := config.Config{}
	app := &cli.App{
		Name:  "shortu",
		Usage: "short your url",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "listen",
				Usage:       "listen port number",
				Aliases:     []string{"l"},
				Value:       8082,
				Destination: &conf.Port,
			},
			&cli.StringFlag{
				Name:        "domain",
				Usage:       "short URL prefix(like https://u.kfd.me)",
				Value:       "https://u.kfd.me",
				Destination: &conf.Prefix,
			},
			&cli.StringFlag{
				Name:        "db-path",
				Aliases:     []string{"p"},
				Usage:       "database path",
				Value:       "short-url.db",
				Destination: &conf.DBPath,
			},
			&cli.StringFlag{
				Name:        "db-type",
				Aliases:     []string{"t"},
				Usage:       "database type: redis or file",
				Value:       "file",
				Destination: &conf.DBType,
			},
			&cli.StringFlag{
				Name:        "redis",
				Aliases:     []string{"r"},
				Usage:       "database path",
				Value:       "localhost:6379/0",
				Destination: &conf.Redis,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "log level: debug",
				Value:       false,
				Destination: &conf.Debug,
			},
		},
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
			return routes.Serve(&conf, &shorter)
		},
	}

	app.Run(os.Args)
}
