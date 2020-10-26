package kubernetes

import (
	"context"
	"io/ioutil"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type cfgKey struct{}
type csKey struct{}

func kubeDir(dir string) string {
	if dir != "" {
		return dir
	}
	return os.Getenv("KUBE_CONFIG_DIR")
}

//GetKubeConfig returns rest config for the give kube config directory
//if directory input is not given its taken from KUBE_CONFIG_DIR env variable
//if that variable is not present incluster config is considered
func GetKubeConfig(dir string) (*rest.Config, error) {
	kubeconfigPath := kubeDir(dir)
	switch {
	case kubeconfigPath != "":
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	case kubeconfigPath == "":
		return rest.InClusterConfig()
	}
	return nil, nil
}

//GetCurrentNamespace returns the namespace where the pod is running
func GetCurrentNamespace() (string, error) {
	dat, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

//GetInClusterKubeClient returns the client set for incluster config
func GetInClusterKubeClient() (*kubernetes.Clientset, error) {
	config, cfgErr := rest.InClusterConfig()
	if cfgErr != nil {
		return nil, cfgErr
	}

	return kubernetes.NewForConfig(config)
}

//WithContext returns the context with
//kubernetes client sets and rest config
func WithContext(ctx context.Context) context.Context {
	restConfig, err := GetKubeConfig("")
	if err != nil {
		panic(err)
	}
	kubernetesCS, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, cfgKey{}, restConfig)
	return context.WithValue(ctx, csKey{}, kubernetesCS)
}

// FromContext returns the config stored in context.
func CfgFromContext(ctx context.Context) *rest.Config {
	if cfg, ok := ctx.Value(cfgKey{}).(*rest.Config); ok {
		return cfg
	}
	return nil
}

// FromContext returns the kubernetes client set stored in context.
func KubernetesCSFromContext(ctx context.Context) *kubernetes.Clientset {
	if cs, ok := ctx.Value(csKey{}).(*kubernetes.Clientset); ok {
		return cs
	}
	return nil
}
