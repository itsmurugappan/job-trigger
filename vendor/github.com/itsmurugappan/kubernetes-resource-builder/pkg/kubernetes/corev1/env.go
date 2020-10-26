package corev1

import (
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
)

func GetEnvfromSecretorCM(envFrom []kubernetes.EnvFrom) []corev1.EnvFromSource {
	if len(envFrom) > 0 && envFrom[0].Name != "" {
		var envs []corev1.EnvFromSource
		for _, env := range envFrom {
			switch env.Type {
			case "Secret":
				envs = append(envs, corev1.EnvFromSource{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: env.Name,
						},
					},
				})
			case "CM":
				envs = append(envs, corev1.EnvFromSource{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: env.Name,
						},
					},
				})
			}
		}
		return envs
	}
	return nil
}

func GetEnvFromHTTPParam(queryParams map[string][]string) []corev1.EnvVar {
	var envVars []corev1.EnvVar
	for k, v := range queryParams {
		if k != "" {
			envVars = append(envVars, corev1.EnvVar{Name: k, Value: strings.Join(v, ";")})
		}
	}
	if len(envVars) > 0 {
		return envVars
	}
	return nil
}
