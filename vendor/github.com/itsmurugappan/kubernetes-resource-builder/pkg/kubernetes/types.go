package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
)

//EnvFrom = env variables
type EnvFrom struct {
	Name string
	Type string
}

//ContainerSpec - kubernetes core/v1/container
type ContainerSpec struct {
	Image             string
	Port              int32
	Name              string
	Resources         []Resource
	Secrets           []corev1.VolumeMount
	ConfigMaps        []corev1.VolumeMount
	EnvVariables      []corev1.EnvVar
	User              int64
	EnvFromSecretorCM []EnvFrom
	Cmd               []string
	ServiceAccount    string
}

//ContainerSpec - kubernetes core/v1/pod
type PodSpec struct {
	Containers []ContainerSpec
}

//JobSpec - kubernetes batch/v1/job
type JobSpec struct {
	Spec PodSpec
	Name string
}

//Resource - container resource constraints
type Resource struct {
	Type string
	CPU  int64
	Mem  int64
}

//KV - generic key/value
type KV struct {
	Key   string
	Value string
}
