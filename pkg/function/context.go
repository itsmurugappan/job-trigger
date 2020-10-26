package function

import (
	"context"
	"encoding/json"
	"os"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
)

type fnCtx struct {
	c context.Context
}

type jobSpecKey struct{}
type nsKey struct{}

func Context(ctx context.Context) *fnCtx {
	return &fnCtx{ctx}
}

func BuildContext() context.Context {
	return withNamespace(
		withSpec(
			kubernetes.WithContext(
				context.Background(),
			),
		),
	)
}

func withSpec(ctx context.Context) context.Context {
	raw, present := os.LookupEnv("JOB_SPEC")
	if !present {
		panic("'JOB_SPEC' not provided as a env variable")
	}
	containerSpec := kubernetes.ContainerSpec{}
	if err := json.Unmarshal([]byte(raw), &containerSpec); err != nil {
		panic(err)
	}
	return context.WithValue(ctx, jobSpecKey{}, &containerSpec)
}

func withNamespace(ctx context.Context) context.Context {
	ns, present := os.LookupEnv("FUNCTION_NAMESPACE")
	if !present {
		ns, _ = kubernetes.GetCurrentNamespace()
	}
	return context.WithValue(ctx, nsKey{}, &ns)
}

func specFromContext(ctx context.Context) *kubernetes.ContainerSpec {
	if spec, ok := ctx.Value(jobSpecKey{}).(*kubernetes.ContainerSpec); ok {
		return spec
	}
	return nil
}

func nsFromContext(ctx context.Context) *string {
	if spec, ok := ctx.Value(nsKey{}).(*string); ok {
		return spec
	}
	return nil
}
