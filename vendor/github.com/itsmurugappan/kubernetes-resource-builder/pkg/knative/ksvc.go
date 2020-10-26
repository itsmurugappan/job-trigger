package knative

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"knative.dev/serving/pkg/client/clientset/versioned"
	typedservingv1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
)

type kClient struct {
	tservingv1 typedservingv1.ServingV1Interface
	ctx        context.Context
}

func Client(c context.Context) *kClient {
	cfg := kubernetes.CfgFromContext(c)
	kcs, _ := versioned.NewForConfig(cfg)

	return &kClient{
		tservingv1: kcs.ServingV1(),
		ctx:        c,
	}
}

//GetKService returns a knative service object for the name and namespaces
func (c kClient) GetKService(ns, name string) (*servingv1.Service, error) {
	return c.tservingv1.Services(ns).Get(c.ctx, name, metav1.GetOptions{})
}
