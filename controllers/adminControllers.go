package controllers

import (
	log "log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminGetEnvConfig(c *gin.Context) {
	log.Info("admin get env config")
	c.JSON(http.StatusOK, gin.H{
		"message": "admin config accessed",
	})
}
