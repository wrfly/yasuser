package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/filter"
	s "github.com/wrfly/yasuser/shortener"
)

const _MaxURLLength = 1 << 11

var urlBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, _MaxURLLength+1)
	},
}

// Serve routes
func Serve(conf config.SrvConfig,
	shortener s.Shortener, filter filter.Filter) error {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	srv := newServer(conf, shortener, filter)

	e := gin.New()

	e.GET("/", srv.handleIndex())
	e.POST("/", srv.handleLongURL())
	e.GET("/:URI", srv.handleURI())

	// go tool pprof
	if conf.Pprof {
		e.GET("/:URI/pprof/", func(c *gin.Context) {
			pprof.Index(c.Writer, c.Request)
		})
		e.GET("/:URI/pprof/:x", func(c *gin.Context) {
			pprof.Index(c.Writer, c.Request)
		})
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: e,
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
