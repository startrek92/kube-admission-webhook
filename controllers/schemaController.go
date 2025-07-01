package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/startrek92/kube-admission-webhook/config"
	"github.com/startrek92/kube-admission-webhook/dao"
	mongomodels "github.com/startrek92/kube-admission-webhook/mongoModels"
	requestmodels "github.com/startrek92/kube-admission-webhook/requestModels"
)

func RequestSchema(c *gin.Context) {
	logrus.Info("Incoming schema check (ping endpoint)")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func MergeEnvConfigs(containers []requestmodels.Container, globalCfg, nsCfg *mongomodels.WorkloadConfig) map[string][]mongomodels.EnvVar {
	finalEnvMap := make(map[string]map[string]string)

	// Step 1: Merge raw values (string -> string)
	for _, container := range containers {
		cn := container.Name
		finalEnvMap[cn] = make(map[string]string)

		logrus.Infof("Processing container: %s", cn)

		// Global config
		if globalCfg != nil {
			if contCfg, ok := globalCfg.Containers[cn]; ok {
				for _, env := range contCfg.Envs {
					finalEnvMap[cn][env.Name] = env.Value
					logrus.Infof("[Global] Set %s=%s", env.Name, env.Value)
				}
			}
		}

		// Namespace config (overwrite)
		if nsCfg != nil {
			if contCfg, ok := nsCfg.Containers[cn]; ok {
				for _, env := range contCfg.Envs {
					finalEnvMap[cn][env.Name] = env.Value
					logrus.Infof("[Namespace] Overwritten %s=%s", env.Name, env.Value)
				}
			}
		}
	}

	// Step 2: Convert to []EnvVar
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
		logrus.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	containers := request.Request.Object.Spec.Template.Spec.Containers
	namespace := request.Request.Namespace
	workloadID := request.Request.Name

	cfg := config.GetConfig()

	// --- Fetch global config ---
	globalCollection, err := cfg.GetCollection("global")
	if err != nil {
		logrus.Errorf("Error getting global collection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal config error"})
		return
	}
	globalCfg, err := dao.GetWorkloadEnv(globalCollection, workloadID)
	if err != nil {
		logrus.Errorf("Error fetching global config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch global config"})
		return
	}

	// --- Fetch namespace config ---
	nsCollection, err := cfg.GetCollection(namespace)
	if err != nil {
		logrus.Warnf("No collection mapping for namespace: %s", namespace)
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown namespace"})
		return
	}
	nsCfg, err := dao.GetWorkloadEnv(nsCollection, workloadID)
	if err != nil {
		logrus.Errorf("Error fetching namespace config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch namespace config"})
		return
	}

	// --- Merge configs ---
	finalEnvs := MergeEnvConfigs(containers, globalCfg, nsCfg)

	logrus.Infof("Final merged envs: %+v", finalEnvs)

	c.JSON(http.StatusOK, gin.H{
		"merged_env_config": finalEnvs, // map[string][]models.EnvVar
	})
}
