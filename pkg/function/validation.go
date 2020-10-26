package function

import (
	"context"

	pkgcorev1 "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes/corev1"
)

func ValidateK8sResource(ctx context.Context) error {
	ns := nsFromContext(ctx)
	spec := specFromContext(ctx)

	if err := pkgcorev1.Client(ctx).CheckSecretMounts(*ns, spec.Secrets); err != nil {
		return err
	}

	if err := pkgcorev1.Client(ctx).CheckCMMounts(*ns, spec.ConfigMaps); err != nil {
		return err
	}

	if err := pkgcorev1.Client(ctx).CheckEnvFromResources(*ns, spec.EnvFromSecretorCM); err != nil {
		return err
	}
	return nil
}
