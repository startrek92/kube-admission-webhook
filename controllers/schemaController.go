package controllers

import (
	"net/http"

	log "log/slog"

	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/config"
	"github.com/startrek92/kube-admission-webhook/dao"
	mongomodels "github.com/startrek92/kube-admission-webhook/mongoModels"
	requestmodels "github.com/startrek92/kube-admission-webhook/requestModels"
	"github.com/startrek92/kube-admission-webhook/services"
)

func RequestSchema(c *gin.Context) {
	log.Info("Incoming schema check (ping endpoint)")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func MergeEnvConfigs(containers []requestmodels.Container, globalCfg, nsCfg *mongomodels.WorkloadConfig) map[string][]mongomodels.EnvVar {
	finalEnvMap := make(map[string]map[string]string)

	for _, container := range containers {
		cn := container.Name
		finalEnvMap[cn] = make(map[string]string)

		log.Info("Processing container", "name", cn)

		if globalCfg != nil {
			if contCfg, ok := globalCfg.Containers[cn]; ok {
				for _, env := range contCfg.Envs {
					finalEnvMap[cn][env.Name] = env.Value
					log.Info("[Global] Set", "container", cn, "name", env.Name, "value", env.Value)
				}
			}
		}

		if nsCfg != nil {
			if contCfg, ok := nsCfg.Containers[cn]; ok {
				for _, env := range contCfg.Envs {
					finalEnvMap[cn][env.Name] = env.Value
					log.Info("[Namespace] Overwritten", "container", cn, "name", env.Name, "value", env.Value)
				}
			}
		}
	}

	finalEnvs := make(map[string][]mongomodels.EnvVar)
	for containerName, envMap := range finalEnvMap {
		var envList []mongomodels.EnvVar
		for k, v := range envMap {
			envList = append(envList, mongomodels.EnvVar{Name: k, Value: v})
		}
		finalEnvs[containerName] = envList
	}

	return finalEnvs
}

func IncomingRequestSchema2(c *gin.Context) {
	var request requestmodels.AdmissionReview
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	containers := request.Request.Object.Spec.Template.Spec.Containers
	namespace := request.Request.Namespace
	workloadID := request.Request.Name

	cfg := config.GetConfig()

	globalCollection, err := cfg.GetCollection("global")
	if err != nil {
		log.Error("Error getting global collection", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal config error"})
		return
	}
	globalCfg, err := dao.GetWorkloadEnv(globalCollection, workloadID)
	if err != nil {
		log.Error("Error fetching global config", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch global config"})
		return
	}

	nsCollection, err := cfg.GetCollection(namespace)
	if err != nil {
		log.Warn("No collection mapping for namespace", "namespace", namespace)
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown namespace"})
		return
	}
	nsCfg, err := dao.GetWorkloadEnv(nsCollection, workloadID)
	if err != nil {
		log.Error("Error fetching namespace config", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch namespace config"})
		return
	}

	finalEnvs := MergeEnvConfigs(containers, globalCfg, nsCfg)

	log.Info("Final merged envs", "data", finalEnvs)

	patchResponse, err := services.BuildPatchResponse(request, finalEnvs)
	if err != nil {
		log.Error("Failed to build patch response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build patch"})
		return
	}

	c.JSON(http.StatusOK, patchResponse)
}
