package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo-contrib/prometheus"

	"github.com/badboyd/dynamo/config"
	"github.com/badboyd/dynamo/pkg/logger"
)

const (
	errFileSizeFmt        = "Please upload file size less than %dMB"
	errInvalidFileTypeFmt = "Please upload these file types: %s"
)

type (
	// Server struct
	Server struct {
		e      *echo.Echo
		cfg    *config.Config
		logger *zap.SugaredLogger
	}

	uploadReq struct {
		FileType string `json:"file_type"`
		FileZize int64  `json:"filze_size"`
	}
)

// H.264 High Profile and VP9 (profile 0)

// New Server
func New(cfg *config.Config) *Server {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	srv := &Server{
		e:      e,
		cfg:    cfg,
		logger: logger.New("server"),
	}

	e.GET("/meta/config", cfg.ServeHTTP)

	e.POST("/video", srv.upload)
	e.GET("/health", srv.checkHealth)

	return srv
}

// Start server. Should run in a separated go-routine
func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.cfg.Server.HTTPPort))
}

// Stop listening for new connections. Finish all running connections
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.e.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) checkHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}

func (s *Server) upload(ctx echo.Context) error {
	fh, err := ctx.FormFile("video")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	upReq := uploadReq{
		FileType: fh.Header.Get("Content-Type"),
		FileZize: fh.Size,
	}

	if err := s.validateReq(ctx, upReq); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, upReq)
}

func (s *Server) validateReq(ctx echo.Context, req uploadReq) error {
	s.logger.Debugf("validate upload req: %+v", req)

	if req.FileZize > s.cfg.Video.MaxSize {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf(errFileSizeFmt, s.cfg.Video.MaxSizeMB))
	}

	if s.cfg.Video.AllowedTypes[strings.TrimPrefix(req.FileType, "video/")] == false {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf(errInvalidFileTypeFmt, s.cfg.Video.AllowedTypesString))
	}

	return nil
}
