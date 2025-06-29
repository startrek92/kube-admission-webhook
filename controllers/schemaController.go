package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/startrek92/kube-admission-webhook/config"
	"github.com/startrek92/kube-admission-webhook/dao"
	"github.com/startrek92/kube-admission-webhook/models"
)

func RequestSchema(c *gin.Context) {
	logrus.Info("Incoming schema check (ping endpoint)")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func IncomingRequestSchema(c *gin.Context) {
	var request models.AdmissionReview
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	cfg := config.GetConfig()

	collectionName, err := cfg.GetCollection(request.Request.Namespace)
	if err != nil {
		logrus.Warnf("No collection mapping found for namespace: %s", request.Request.Namespace)
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown namespace"})
		return
	}

	logrus.Infof("Resolved collection: %s for namespace: %s", collectionName, request.Request.Namespace)

	workloadID := request.Request.Name

	envData, err := dao.GetWorkloadEnv(collectionName, workloadID)
	if err != nil {
		logrus.Errorf("Failed to fetch workload env: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch workload config"})
		return
	}

	if envData == nil {
		logrus.Infof("No env config found for workload: %s in %s", workloadID, collectionName)
		c.JSON(http.StatusNotFound, gin.H{"message": "no env config found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workload_env": envData,
	})
}
