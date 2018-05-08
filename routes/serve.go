package routes

import (
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/short-url/config"
	"github.com/wrfly/short-url/handler"
)

const MAX_URL_LENGTH = 1e3

var urlBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, MAX_URL_LENGTH+1)
	},
}

// Serve routes
func Serve(conf *config.Config, shorter *handler.Shorter) error {
	example := fmt.Sprintf("curl %s -d \"%s\"",
		conf.Prefix, "http://longlonglong.com/long/long/long?a=1&b=2")

	srv := gin.Default()

	srv.GET("/:s", func(c *gin.Context) {
		location := shorter.Long(c.Param("s"))
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
		buf := urlBufferPool.Get().([]byte)
		defer urlBufferPool.Put(buf)
		n, err := c.Request.Body.Read(buf)
		if err != io.EOF && err != nil {
			c.String(500, fmt.Sprintf("error: %s\n", err))
			return
		}
		if n > MAX_URL_LENGTH {
			c.String(400, fmt.Sprintln("Bad request, URL too long"))
			return
		}

		longURL := fmt.Sprintf("%s", buf[:n])
		short := shorter.Short(longURL)
		shortURL := fmt.Sprintf("%s/%s", conf.Prefix, short)
		logrus.Debugf("shorten URL: [ %s ] -> [ %s ]",
			longURL, shortURL)
		c.String(200, fmt.Sprintln(shortURL))
	})

	port := fmt.Sprintf(":%d", conf.Port)
	logrus.Infof("Service running at [ %s ], with prefix [ %s ]",
		port, conf.Prefix)

	return srv.Run(port)
}
