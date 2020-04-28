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
	"github.com/labstack/gommon/random"

	"github.com/labstack/echo-contrib/prometheus"

	"github.com/badboyd/dynamo/config"
	"github.com/badboyd/dynamo/pkg/logger"
	"github.com/badboyd/dynamo/pkg/storage"
	"github.com/badboyd/dynamo/pkg/storage/gcs"
	"github.com/badboyd/dynamo/pkg/storage/s3"
)

const (
	errFileSizeFmt        = "Please upload file size less than %dMB"
	errInvalidFileTypeFmt = "Please upload these file types: %s"
)

type (
	// Server struct
	Server struct {
		e       *echo.Echo
		cfg     *config.Config
		logger  *zap.SugaredLogger
		storage storage.Storage
	}

	uploadReq struct {
		FileType string `json:"file_type"`
		FileZize int64  `json:"filze_size"`
		FileName string `json:"file_name"`
		FileURL  string `json:"file_url"`
	}
)

// H.264 High Profile and VP9 (profile 0)

// New Server
func New(cfg *config.Config) *Server {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	p := prometheus.NewPrometheus("dynamo", nil)
	p.Use(e)

	srv := &Server{
		e:      e,
		cfg:    cfg,
		logger: logger.New("server"),
	}

	var err error
	if cfg.GCS.Enable {
		srv.storage, err = gcs.New(cfg.GCS.Bucket)
	}
	if err != nil {
		panic(err)
	}

	if srv.storage == nil && cfg.S3.Enable {
		srv.storage, err = s3.New(cfg.S3.Bucket)
	}
	if err != nil {
		panic(err)
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

	if err := s.storage.Close(); err != nil {
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

	fileName := fmt.Sprintf("raw/%s", random.New().String(10))
	upReq := uploadReq{
		FileType: fh.Header.Get("Content-Type"),
		FileZize: fh.Size,
		FileName: fileName,
		FileURL:  fmt.Sprintf("https://cdn.chotot.org/%s", fileName),
	}

	if err := s.validateReq(ctx, upReq); err != nil {
		return err
	}

	f, err := fh.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer f.Close()

	if err := s.storage.Write(context.Background(), f, fileName, true); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
