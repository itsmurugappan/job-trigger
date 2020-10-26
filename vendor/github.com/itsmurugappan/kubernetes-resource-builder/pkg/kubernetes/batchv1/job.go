package batchv1

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedbatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"

	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/ptr"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes/corev1"
	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/transform"
)

type JobSpecOption func(*batchv1.Job)

type batchClient struct {
	tbatchv1 typedbatchv1.BatchV1Interface
	ctx      context.Context
}

func Client(c context.Context) *batchClient {
	cs := kubernetes.KubernetesCSFromContext(c)
	return &batchClient{
		tbatchv1: cs.BatchV1(),
		ctx:      c,
	}
}

func (c *batchClient) CreateJob(ns string, job *batchv1.Job, watch bool) (*batchv1.Job, error) {
	return c.tbatchv1.Jobs(ns).Create(c.ctx, job, metav1.CreateOptions{})
}

func (c *batchClient) GetJobStatus(ns, jobName string) (*batchv1.JobStatus, error) {
	job, err := c.tbatchv1.Jobs(ns).Get(c.ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &job.Status, nil
}

func GetJob(name string, options ...JobSpecOption) batchv1.Job {
	jobSpec := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		}}

	for _, fn := range options {
		fn(&jobSpec)
	}
	return jobSpec
}

func WithPodSpecOptions(podSpec kubernetes.PodSpec, options ...corev1.PodSpecOption) JobSpecOption {
	return func(job *batchv1.Job) {
		job.Spec.Template.Spec = corev1.GetPodSpec(podSpec, options...)
	}
}

func WithTTL(ttl int32) JobSpecOption {
	return func(job *batchv1.Job) {
		if ttl > int32(0) {
			job.Spec.TTLSecondsAfterFinished = ptr.Int32(ttl)
		}
	}
}

func WithParallelism(instances int32) JobSpecOption {
	return func(job *batchv1.Job) {
		if instances > int32(0) {
			job.Spec.Parallelism = ptr.Int32(instances)
		}
	}
}

func WithBackoffLimit(backoffLimit int32) JobSpecOption {
	return func(job *batchv1.Job) {
		job.Spec.BackoffLimit = ptr.Int32(backoffLimit)
	}
}

func WithAnnotations(inAnnotations []kubernetes.KV) JobSpecOption {
	return func(job *batchv1.Job) {
		job.Spec.Template.ObjectMeta.Annotations = transform.GetStringMap(inAnnotations, nil)
	}
}

func WithLabels(inLabels []kubernetes.KV) JobSpecOption {
	return func(job *batchv1.Job) {
		job.Spec.Template.ObjectMeta.Labels = transform.GetStringMap(inLabels, nil)
	}
}

func WithOwnerReference(obj kmeta.OwnerRefable) JobSpecOption {
	return func(job *batchv1.Job) {
		ownerRef := metav1.NewControllerRef(obj.GetObjectMeta(), obj.GetGroupVersionKind())
		job.ObjectMeta.OwnerReferences = []metav1.OwnerReference{*ownerRef}
	}
}
