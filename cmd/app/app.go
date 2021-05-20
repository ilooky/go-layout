package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ilooky/go-layout/pkg/config"
	"github.com/ilooky/go-layout/pkg/database"
	"github.com/ilooky/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type app struct {
	server *http.Server
	errCh  chan error
	conf   *config.Config
	port   int
}

func newApp(conf *config.Config, api func(ctx *gin.Engine)) *app {
	port, _ := strconv.Atoi(conf.Port)
	if conf.Log.Release {
		gin.SetMode("release")
	}
	h := gin.New()
	h.RemoveExtraSlash = true
	h.RedirectFixedPath = true
	h.Use(gin.Recovery(), logMiddleware())
	h.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "SUCCESS")
	})
	api(h)
	a := app{
		server: &http.Server{
			Addr:    ":" + conf.Port,
			Handler: h,
		},
		errCh: make(chan error),
		conf:  conf,
		port:  port,
	}
	return &a
}

func (a app) start() {
	go func() {
		if err := config.Client.Register(a.conf.Name, a.conf.Host, a.port); err != nil {
			a.errCh <- err
		}
	}()
	go func() {
		a.errCh <- a.server.ListenAndServe()
	}()
}

func (a app) shutdown() {
	_ = config.Client.UnRegister(a.conf.Name, a.conf.Host, a.port)
	_ = a.server.Close()
}

func (a app) error() chan error {
	return a.errCh
}

func Run(serverName string, server func(g *gin.Engine), inits ...func(cnf *config.Config)) error {
	cloud, err := config.InitCloud()
	if err != nil {
		return err
	}
	conf, err := cloud.ReadConfig(serverName)
	if err != nil {
		return err
	}
	logger.InitLogger(logger.Config{
		Level: conf.Log.Level,
		Path:  conf.Log.Path,
	})
	logger.InfoKV("Read Config", conf.Name, conf)
	if db, err := database.InitOrm(conf.Mysql); err != nil {
		logger.Panic(err)
		return err
	} else {
		defer db.Close()
	}
	app := newApp(conf, server)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	app.start()
	for _, init := range inits {
		init(conf)
	}
	logger.Infof("服务已启动...监听[%s]端口", conf.Port)
	defer app.shutdown()
	select {
	case err := <-app.errCh:
		return err
	case <-quit:
		return nil
	}
}

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		start := time.Now()
		c.Next()
		if path != "/health" {
			if raw != "" {
				path = path + "?" + raw
			}
			method := c.Request.Method
			statusCode := c.Writer.Status()
			errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()
			msg := struct {
				Method     string
				StatusCode int
				ERROR      string
				Latency    time.Duration
			}{
				Method:     method,
				StatusCode: statusCode,
				ERROR:      errMsg,
				Latency:    time.Now().Sub(start),
			}
			logger.InfoKv("Request", zap.Any(path, msg))
		}
	}
}
