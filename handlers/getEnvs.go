package handlers

import "github.com/startrek92/kube-admission-webhook/models"


func getEnvVars(request models.AdmissionReview) {
	var namespace = request.Request.Namespace;
	var deploymentName = request.Request.Name;
	var operation = request.Request.Operation;
}