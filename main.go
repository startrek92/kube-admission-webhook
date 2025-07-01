package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/config"
	"github.com/startrek92/kube-admission-webhook/controllers"
	"github.com/startrek92/kube-admission-webhook/db"
	"github.com/startrek92/kube-admission-webhook/logger"
)

func main() {
	config.InitConfig("./config/config.toml")
	cfg := config.GetConfig()

	// Initialize structured logger
	logger.InitLogger(cfg)
	logger.Log.Info("Server starting...")

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(logger.GinLogrusMiddleware())

	router.GET("", func(c *gin.Context) {
		logger.Log.Info("Health check route hit")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/update", controllers.IncomingRequestSchema2)

	logger.Log.Infof("Connecting to MongoDB at %s", cfg.BuildMongoURI())
	db.Connect(cfg.BuildMongoURI())

	logger.Log.Infof("Running HTTPS server on %s", serverAddr)
	err := router.RunTLS(serverAddr, cfg.Server.CertFile, cfg.Server.KeyFile)
	if err != nil {
		logger.Log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}
