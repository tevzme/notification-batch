package routes

import (
	"net/http"

	"notification_batch/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

// Init initializes the Gin router and sets up routes.
func Init(router *gin.Engine, sch *gocron.Scheduler, cfgMap map[string]*config.Config) {
	setupRoutes(router)
}

func setupRoutes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "connected!"})
	})
}
