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

	srv := server{
		domain: conf.Domain,
		stener: shortener,
	}
	srv.init()

	engine := gin.New()
	engine.GET("/", srv.handleIndex())
	engine.POST("/", srv.handleLongURL())
	engine.GET("/:URI", srv.handleURI())

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: engine,
	}

	errChan := make(chan error)
	go func() {
		errChan <- httpServer.ListenAndServe()
	}()
	logrus.Infof("Server running at [ http://0.0.0.0:%d ], with domain [ %s ]",
		conf.Port, conf.Domain)

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
