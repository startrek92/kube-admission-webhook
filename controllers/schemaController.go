package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/models"
)

func RequestSchema(c *gin.Context) {
	fmt.Println("incoming schema received")
	c.JSON(200, gin.H {
		"status": "ok",
	})
}

func IncomingRequestSchema(c *gin.Context) {
	var request models.AdmissionReview;
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("invalid request ", err);
		c.JSON(400, gin.H{"error": "invalid request body"});
		return;
	}
}

