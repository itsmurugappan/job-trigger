package corev1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	types "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
)

const (
	//CM_MISSING - error message to indicate CM is missing
	CM_MISSING = "cm %s is not found in the namespace %s, create cm before referring in the function"
	//SECRET_MISSING - error message to indicate Secret is missing
	SECRET_MISSING = "secret %s is not found in the namespace %s, create secret before referring in the function"
	//INVALID_ENV_FROM_TYPE - error message to indicate invalid type to mount as env
	INVALID_ENV_FROM_TYPE = "Provide a valid EnvFrom type. Should be 'CM' or 'Secret'"
)

func (c *coreClient) CheckIfCMExist(nsName string, cm string) bool {
	getOpts := metav1.GetOptions{TypeMeta: metav1.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"}}

	if _, err := c.tcorev1.ConfigMaps(nsName).Get(c.ctx, cm, getOpts); err != nil {
		return false
	}

	return true
}

func (c *coreClient) CheckIfSecretExist(nsName string, secret string) bool {
	getOpts := metav1.GetOptions{TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"}}

	if _, err := c.tcorev1.Secrets(nsName).Get(c.ctx, secret, getOpts); err != nil {
		return false
	}

	return true
}

func (c *coreClient) CheckSecretMounts(ns string, secrets []corev1.VolumeMount) error {
	// check secret
	for _, secret := range secrets {
		if secret.Name == "" {
			continue
		}
		if !c.CheckIfSecretExist(ns, secret.Name) {
			return fmt.Errorf(SECRET_MISSING, secret.Name, ns)
		}
	}
	return nil
}

func (c *coreClient) CheckCMMounts(ns string, cms []corev1.VolumeMount) error {
	// check cm
	for _, cm := range cms {
		if cm.Name == "" {
			continue
		}
		if !c.CheckIfCMExist(ns, cm.Name) {
			return fmt.Errorf(CM_MISSING, cm.Name, ns)
		}
	}
	return nil
}

func (c *coreClient) CheckEnvFromResources(ns string, envsFrom []types.EnvFrom) error {
	// check env from
	for _, env := range envsFrom {
		if env.Name == "" {
			continue
		}
		switch env.Type {
		case "CM":
			if !c.CheckIfCMExist(ns, env.Name) {
				return fmt.Errorf(CM_MISSING, env.Name, ns)
			}
		case "Secret":
			if !c.CheckIfSecretExist(ns, env.Name) {
				return fmt.Errorf(SECRET_MISSING, env.Name, ns)
			}
		default:
			return fmt.Errorf(INVALID_ENV_FROM_TYPE)
		}
	}
	return nil
}

func (c *coreClient) GetSecrets(nsName, secret string) (*corev1.Secret, error) {
	getOpts := metav1.GetOptions{TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"}}

	return c.tcorev1.Secrets(nsName).Get(c.ctx, secret, getOpts)
}

func (c *coreClient) GetSAToken(nsName, saName string) (string, error) {
	listOpts := metav1.ListOptions{FieldSelector: "type=kubernetes.io/service-account-token"}

	secretList, err := c.tcorev1.Secrets(nsName).List(c.ctx, listOpts)
	if err != nil {
		return "", err
	}
	for _, secret := range secretList.Items {
		if secret.ObjectMeta.Annotations["kubernetes.io/service-account.name"] == saName {
			return string(secret.Data["token"]), nil
		}
	}
	return "", fmt.Errorf("No secret for sa %s", saName)
}
