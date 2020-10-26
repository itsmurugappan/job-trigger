package corev1

import (
	corev1 "k8s.io/api/core/v1"
)

func GetVolumeSources(configmaps []corev1.VolumeMount, secrets []corev1.VolumeMount) []corev1.Volume {
	var volList []corev1.Volume

	//iterate cm & secret
	for _, cm := range configmaps {
		if cm.Name != "" {
			vol := corev1.Volume{
				Name: cm.Name, VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: cm.Name,
						},
					},
				},
			}
			volList = append(volList, vol)
		}
	}
	for _, secret := range secrets {
		if secret.Name != "" {
			vol := corev1.Volume{
				Name: secret.Name, VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: secret.Name,
					},
				},
			}
			volList = append(volList, vol)
		}
	}
	return volList
}

func GetVolumeMounts(configmaps []corev1.VolumeMount, secrets []corev1.VolumeMount) []corev1.VolumeMount {
	var mountList []corev1.VolumeMount

	if len(configmaps) > 0 && configmaps[0].Name != "" {
		mountList = append(mountList, configmaps...)
	}
	if len(secrets) > 0 && secrets[0].Name != "" {
		mountList = append(mountList, secrets...)
	}
	return mountList
}
