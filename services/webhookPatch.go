package services

import (
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	mongomodels "github.com/startrek92/kube-admission-webhook/mongoModels"
	requestmodels "github.com/startrek92/kube-admission-webhook/requestModels"
)

type jsonPatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func BuildPatchResponse(
	req requestmodels.AdmissionReview,
	merged map[string][]mongomodels.EnvVar,
) (*admissionv1.AdmissionReview, error) {
	containers := req.Request.Object.Spec.Template.Spec.Containers
	var patches []jsonPatchOp

	for i, c := range containers {
		envs, ok := merged[c.Name]
		if !ok || len(envs) == 0 {
			continue
		}

		patches = append(patches, jsonPatchOp{
			Op:    "replace",
			Path:  fmt.Sprintf("/spec/template/spec/containers/%d/env", i),
			Value: envs,
		})
	}

	if len(patches) == 0 {
		return &admissionv1.AdmissionReview{
			Response: &admissionv1.AdmissionResponse{
				UID:     types.UID(req.Request.UID),
				Allowed: true,
			},
		}, nil
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patch: %w", err)
	}

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:       types.UID(req.Request.UID),
			Allowed:   true,
			PatchType: ptr(admissionv1.PatchTypeJSONPatch),
			Patch:     patchBytes,
		},
	}, nil
}

func ptr[T any](v T) *T { return &v }
