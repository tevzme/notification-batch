package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"notification_batch/internal/config"
	"notification_batch/internal/logger"
	"notification_batch/internal/routes"
	"notification_batch/internal/scheduler"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func main() {
	// Load Configuration
	cfgMap := config.LoadConfig(config.GetEnv())
	defaultCfg, ok := cfgMap["default"]
	if !ok {
		log.Fatalf("Failed to load default configuration")
		return
	}

	// Initialize Logger
	err := logger.InitLogger(defaultCfg.LogPath)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		return
	}
	defer logger.AppLogger.Sync() // Flush buffered log entries

	logger.AppLogger.Sugar().Info("Application is starting with Gin Framework...")

	// Initialize Gin Router
	router = gin.Default()

	// Initialize Scheduler
	scheduler.InitScheduler(cfgMap)

	// Initialize Application and Setup Routes
	routes.Init(router, nil, cfgMap)

	// Start the Scheduler
	scheduler.StartScheduler()

	// Start the Gin HTTP Server
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080" // Default port
		}
		logger.AppLogger.Sugar().Infof("Gin server listening on port %s", port)
		if err := router.Run(":" + port); err != nil && err != http.ErrServerClosed {
			logger.AppLogger.Sugar().Fatalf("Failed to start Gin server: %v", err)
		}
	}()

	// Wait for termination signals (Ctrl+C or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.AppLogger.Info("Application is shutting down...")

	// Stop the Scheduler
	scheduler.StopScheduler()

	logger.AppLogger.Info("Gin server shutting down...")

	logger.AppLogger.Info("Application has stopped.")
}
