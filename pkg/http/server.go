package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-microservice/internal/config"
	"go-microservice/internal/middleware"
	"go-microservice/internal/routes"
	"go-microservice/pkg/logger"
)

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	cfg := config.GetConfig().Server

	setGinMode(cfg.Environment)

	engine := gin.New()
	server := &Server{
		engine: engine,
	}
	server.registerMiddlewares()

	return server
}

func setGinMode(env string) {
	switch env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	logger.Info("Gin mode set", zap.String("mode", gin.Mode()))
}

func (s *Server) registerMiddlewares() {
	cfg := config.GetConfig().Server

	s.engine.Use(gin.Recovery())
	s.engine.Use(middleware.LoggingMiddleware())
	s.engine.Use(middleware.GzipMiddleware())
	logger.Info("Default middlewares initialized: Recovery, Logging, Gzip")

	if cfg.AllowOrigins != "" {
		s.engine.Use(middleware.CORSMiddleware(cfg.AllowOrigins))
		logger.Info("CORS middleware applied", zap.String("allowed_origins", cfg.AllowOrigins))
	} else {
		logger.Info("CORS middleware skipped, no allowed origins configured")
	}
}

func (s *Server) RegisterRoutes() error {
	routes.RegisterRoutes(s.engine)
	return nil
}

func (s *Server) Start(ctx context.Context) error {
	cfg := config.GetConfig().Server

	httpSrv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: s.engine,
	}

	go func() {
		logger.Info("HTTP server is listening...",
			zap.String("port", cfg.Port))
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server failed unexpectedly", zap.Error(err))
		}

		logger.Info("HTTP server stopped.")
	}()

	<-ctx.Done()
	logger.Info("Shutdown signal received, stopping HTTP server...")

	if err := httpSrv.Shutdown(context.Background()); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("HTTP server shut down gracefully")
	return nil
}
