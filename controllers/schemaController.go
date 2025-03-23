package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RequestSchema(c *gin.Context) {
	fmt.Println("incoming schema received")
	c.JSON(200, gin.H {
		"status": "ok",
	})
}