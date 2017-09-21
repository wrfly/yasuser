package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/wrfly/short-url/config"
	"github.com/wrfly/short-url/handler"
	"github.com/wrfly/short-url/handler/db"
	"gopkg.in/urfave/cli.v2"
)

func srv(conf *config.Config, shorter *handler.Shorter) error {
	example := "curl https://u.kfd.me -d \"http://longlonglong.com/long/long/long?a=1&b=2\""
	srv := gin.Default()
	srv.GET("/:url", func(c *gin.Context) {
		location := shorter.Long(c.Param("url"))
		if location == "" {
			c.String(404, example)
			return
		}
		c.Redirect(302, location)
	})
	srv.GET("/", func(c *gin.Context) {
		c.String(200, example)
	})
	srv.POST("/", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logrus.Error(err)
			c.String(500, err.Error())
		}
		short := shorter.Short(string(b))
		shortURL := fmt.Sprintf("%s/%s\n", conf.Prefix, short)
		c.String(200, shortURL)
	})
	port := fmt.Sprintf(":%d", conf.Port)
	return srv.Run(port)
}

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
				Name:        "prefix",
				Usage:       "domain prefix",
				Aliases:     []string{"p"},
				Value:       "https://u.kfd.me",
				Destination: &conf.Prefix,
			},
			&cli.StringFlag{
				Name:        "db-path",
				Aliases:     []string{"db"},
				Usage:       "database path",
				Value:       "short-url.db",
				Destination: &conf.DBPath,
			},
		},
		Action: func(c *cli.Context) error {
			blotDB, err := db.NewDB(conf.DBPath)
			if err != nil {
				logrus.Fatal(err)
			}
			shorter := handler.Shorter{
				DB: blotDB,
			}
			return srv(&conf, &shorter)
		},
	}

	app.Run(os.Args)
}
