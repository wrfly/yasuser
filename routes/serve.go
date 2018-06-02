package routes

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/config"
	stner "github.com/wrfly/yasuser/shortener"
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
	logrus.Infof("Service starting at [ http://0.0.0.0%s ], with prefix [ %s ]",
		port, conf.Prefix)

	srv := gin.New()

	srv.GET("/", handleIndex(conf.Prefix))
	srv.GET("/:s", handleShortURL(shortener))
	srv.POST("/", handleLongURL(conf.Prefix, shortener))

	return srv.Run(port)
}
