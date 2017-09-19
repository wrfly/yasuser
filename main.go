package main

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/wrfly/short-url/config"
	"github.com/wrfly/short-url/handler"
	"github.com/wrfly/short-url/handler/db"
	"gopkg.in/urfave/cli.v2"
)

func srv(conf config.Config, shorter *handler.Shorter) error {
	srv := gin.Default()
	srv.GET("/:url", func(c *gin.Context) {
		location := shorter.Long(c.Param("url"))
		if location == "" {
			c.String(404, "URL Not Found")
			return
		}
		c.Redirect(302, location)
	})
	srv.POST("/", func(c *gin.Context) {
		short := shorter.Short(c.PostForm("url"))
		shortURL := path.Join(conf.Prefix, short)
		c.String(200, shortURL)
	})
	port := fmt.Sprintf(":%d", conf.Port)
	return srv.Run(port)
}

func main() {
	conf := config.Config{}
	blotDB, err := db.NewDB("/tmp/test.db")
	if err != nil {
		logrus.Fatal(err)
	}
	shorter := handler.Shorter{
		DB: blotDB,
	}
	app := &cli.App{
		Name:  "shortu",
		Usage: "short your url",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "listen",
				Aliases:     []string{"l"},
				Value:       8082,
				Destination: &conf.Port,
			},
			&cli.StringFlag{
				Name:        "prefix",
				Aliases:     []string{"p"},
				Value:       "https://url.kfd.me",
				Destination: &conf.Prefix,
			},
		},
		Action: func(c *cli.Context) error {
			return srv(conf, &shorter)
		},
	}

	app.Run(os.Args)
}
