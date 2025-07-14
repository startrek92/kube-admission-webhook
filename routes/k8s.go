package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/controllers"
)

func AddK8sRoutes(r *gin.Engine) {
	r.POST("/update", controllers.IncomingRequestSchema2)
}
