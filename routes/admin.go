package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/config"
	"github.com/startrek92/kube-admission-webhook/middleware"
)

func AddAdminRoutes(r *gin.Engine, config *config.Config) {
	group := r.Group("admin", middleware.AdminAuthMiddleware(*config))

	group.GET("/config", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "success",
			"msg":    "test route",
		})
	})

}
