package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

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
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	engine := gin.New()
	engine.GET("/", handleIndex(conf.Prefix))
	engine.GET("/:s", handleShortURL(shortener))
	engine.POST("/", handleLongURL(conf.Prefix, shortener))

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: engine,
	}

	errChan := make(chan error)
	go func() {
		errChan <- httpServer.ListenAndServe()
	}()
	logrus.Infof("Server running at [ http://0.0.0.0:%s ], with prefix [ %s ]",
		conf.Port, conf.Prefix)

	select {
	case <-sigChan:
		logrus.Info("Shuting down the server")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		logrus.Info("Server shutdown")
		return err
	case err := <-errChan:
		return err
	}
}
