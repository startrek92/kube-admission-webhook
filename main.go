package main

import (
	"fmt"
	"net/http"

	log "log/slog"

	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/config"
	"github.com/startrek92/kube-admission-webhook/controllers"
	"github.com/startrek92/kube-admission-webhook/db"
	"github.com/startrek92/kube-admission-webhook/logger"
)

func main() {
	config.InitConfig("./config/config.toml")
	cfg := config.GetConfig()

	if err := logger.InitLogger(cfg); err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}
	log.Info("Server starting...")

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.GinSlogMiddleware())

	router.GET("", func(c *gin.Context) {
		log.Info("Health check route hit")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/update", controllers.IncomingRequestSchema2)

	db.Connect(cfg.BuildMongoURI())

	log.Info("Running HTTPS server", "address", serverAddr)
	err := router.RunTLS(serverAddr, cfg.Server.CertFile, cfg.Server.KeyFile)
	if err != nil {
		log.Error("Failed to start HTTPS server", "error", err)
	}
}
