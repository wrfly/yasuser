package routes

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/short-url/config"
	"github.com/wrfly/short-url/handler"
)

func Serve(conf *config.Config, shorter *handler.Shorter) error {
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
	logrus.Infof("Service running at [ %s ], with prefix [ %s ]",
		port, conf.Prefix)

	return srv.Run(port)
}
