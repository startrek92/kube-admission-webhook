package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/controllers"
	"github.com/startrek92/kube-admission-webhook/db"
	"github.com/startrek92/kube-admission-webhook/initializers"
)

func init() {
	initializers.LoadEnvVariables()
}

func LogRequestBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Log the request body
		log.Println("Request Body:", string(bodyBytes))

		// Restore the request body (since it's read once)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Continue to the next handler
		c.Next()
	}
}

func main() {
	certFile := "./certs/webhook-tls.crt"
	keyFile := "./certs/webhook-tls.key"
	serverAddr := "192.168.1.10:5555"

	router := gin.Default()
	// router.Use(LogRequestBodyMiddleware())
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/update", controllers.RequestSchema)
	db.Connect("mongoDbConnectionString");
	router.RunTLS(serverAddr, certFile, keyFile)
}
