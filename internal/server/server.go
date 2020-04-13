package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/badboyd/dynamo/config"
	"github.com/badboyd/dynamo/pkg/storage"
)

type (
	Server struct {
		e          *echo.Echo
		gcsStorage *storage.GCS
		cfg        *config.Config
	}
)

// H.264 High Profile and VP9 (profile 0)

func New(cfg *config.Config) *Server {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	s := &Server{}

	e.GET("/health", s.checkHealth)
	e.GET("/meta/config", cfg.ServeHTTP)
	e.POST("/video", s.uploadVideo)

	return &Server{
		e:          e,
		cfg:        cfg,
		gcsStorage: storage.NewGCS(),
	}
}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.cfg.Server.HTTPPort))
}

func (s *Server) Close() {
	s.e.Close()
}

func (s *Server) checkHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}

func (s *Server) uploadVideo(ctx echo.Context) error {
	// TODO: Implement upload video
	return ctx.JSON(http.StatusOK, "OK")
}
