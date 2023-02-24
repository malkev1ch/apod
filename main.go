package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	genv1 "github.com/malkev1ch/apod/gen/v1"
	"github.com/malkev1ch/apod/internal/config"
	handlerv1 "github.com/malkev1ch/apod/internal/handler/v1"
	middlewarev1 "github.com/malkev1ch/apod/internal/handler/v1/middleware"
	"github.com/malkev1ch/apod/internal/repository"
	"github.com/malkev1ch/apod/internal/service"
	"github.com/malkev1ch/apod/pkg/logger"
	"github.com/malkev1ch/apod/pkg/postgres"
	"github.com/malkev1ch/apod/pkg/s3"
)

var version string

const httpClientConnectionTimeout = time.Second * 20

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed create new app config: %v", err)
	}

	l := logger.New(logger.NewConfig(cfg.LogLevel, cfg.DevMode))
	ctx = logger.ContextWithLogger(ctx, l)

	postgresPool, err := postgres.NewPool(cfg.PGConnectionString)
	if err != nil {
		l.Fatalw(
			"failed to connect to postgres",
			"error", err,
			"url", cfg.PGConnectionString,
		)
	}

	httpClient := http.Client{
		Timeout: httpClientConnectionTimeout,
	}

	s3Client, err := s3.NewS3Client(cfg.S3Address, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Secured)
	if err != nil {
		l.Fatalw(
			"failed to connect to s3",
			"error", err,
			"address", cfg.S3Address,
			"accessKey", cfg.S3AccessKey,
			"secretKey", cfg.S3SecretKey,
			"secured", cfg.S3Secured,
		)
	}

	middlewareManager := middlewarev1.NewManager(l)

	nasaAPIClient := repository.NewNasaAPI(&httpClient, cfg.NASAAPIKey, cfg.NASAPIAddress)
	pictureRepository := repository.NewPicture(postgresPool)
	fileRepository := repository.NewFileRepository(s3Client, cfg.S3Endpoint)

	pictureService := service.NewPicture(nasaAPIClient, pictureRepository, fileRepository, cfg.S3Bucket)

	pictureHandler := handlerv1.NewPicture(pictureService)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.HideBanner = true
	e.HidePort = true

	e.Use(middlewareManager.InjectLogger, middlewareManager.RequestLogger, middlewareManager.BusinessError)

	genv1.RegisterHandlers(e, pictureHandler)

	routes := e.Routes()
	for _, v := range routes {
		l.Infof("registered route %v with method %v", v.Path, v.Method)
	}

	go func() {
		err = e.Start(cfg.ServerListenAddress)
		if err != nil {
			l.Fatalw(
				"failed to start echo server",
				"error", err,
				"addr", cfg.ServerListenAddress,
			)
		}
	}()

	l.Infow("application successfully started", "addr", cfg.ServerListenAddress, "version", version)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = e.Shutdown(ctx)
	if err != nil {
		l.Fatalw(
			"failed to shutdown echo server",
			"error", err,
		)
	}
}
