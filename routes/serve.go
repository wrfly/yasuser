package routes

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/short-url/config"
	stner "github.com/wrfly/short-url/shortener"
)

const MAX_URL_LENGTH = 1e3

var urlBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, MAX_URL_LENGTH+1)
	},
}

// Serve routes
func Serve(conf config.SrvConfig, shortener stner.Shortener) error {
	port := fmt.Sprintf(":%d", conf.Port)
	logrus.Infof("Service starting at [ %s ], with prefix [ %s ]",
		port, conf.Prefix)

	srv := gin.New()

	srv.GET("/", handleIndex(conf.Prefix))
	srv.GET("/:s", handleShortURL(shortener))
	srv.POST("/", handleLongURL(conf.Prefix, shortener))

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
