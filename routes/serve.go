package routes

import (
	"fmt"
	"io"
	"net/url"
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
		realURL := shorter.Long(c.Param("s"))
		if realURL == "" {
			c.String(404, fmt.Sprintln("not found"))
			return
		}
		c.Redirect(302, realURL)
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
		if invalidURL(longURL) {
			c.String(400, fmt.Sprintln("invalid URL"))
			return
		}

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

func invalidURL(URL string) bool {
	logrus.Debugf("get url: %s", URL)
	u, err := url.Parse(URL)
	if err != nil {
		return true
	}

	switch u.Scheme {
	case "":
		return true
	case "http":
	case "https":
	case "ftp":
	case "tcp":
	default:
		return true
	}

	return false
}
