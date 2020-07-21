package function

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedbatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	types "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
	pkgbatchv1 "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes/batchv1"
	pkgcorev1 "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes/corev1"
)

//CreateJob construct a job and creates in the same namespace as the knative service
func CreateJob(spec types.ContainerSpec, ns string, queryParams map[string][]string,
	typedcorev1 typedcorev1.CoreV1Interface,
	typedbatchv1 typedbatchv1.BatchV1Interface) error {
	if err := pkgcorev1.CheckSecretMounts(ns, typedcorev1, spec.Secrets); err != nil {
		return err
	}

	if err := pkgcorev1.CheckCMMounts(ns, typedcorev1, spec.ConfigMaps); err != nil {
		return err
	}

	if err := pkgcorev1.CheckEnvFromResources(ns, typedcorev1, spec.EnvFromSecretorCM); err != nil {
		return err
	}
	labels := getLabelsFromQP(queryParams)
	// delete old jobs
	cleanOldJobs(spec.Name, ns, queryParams["history"], typedbatchv1, typedcorev1)
	delete(queryParams, "history")
	delete(queryParams, "labels")

	jobName := fmt.Sprintf("%s-%s", spec.Name, uuid.New())

	jobSpec := pkgbatchv1.GetJob(types.JobSpec{Name: jobName},
		pkgbatchv1.WithTTL(int32(100)),
		pkgbatchv1.WithBackoffLimit(int32(1)),
		pkgbatchv1.WithLabels(labels),
		pkgbatchv1.WithAnnotations([]types.KV{{"sidecar.istio.io/inject", "false"}}),
		pkgbatchv1.WithPodSpecOptions(types.PodSpec{},
			pkgcorev1.WithRestartPolicy("Never"),
			pkgcorev1.WithVolumes([]types.ContainerSpec{{Secrets: spec.Secrets, ConfigMaps: spec.ConfigMaps}}),
			pkgcorev1.WithContainerOptions(types.ContainerSpec{Image: spec.Image},
				pkgcorev1.WithEnv(spec.EnvVariables),
				pkgcorev1.WithEnv(pkgcorev1.GetEnvFromHTTPParam(queryParams)),
				pkgcorev1.WithEnvFromSecretorCM(spec.EnvFromSecretorCM),
				pkgcorev1.WithVolumeMounts(spec.ConfigMaps, spec.Secrets),
				pkgcorev1.WithSecurityContext(spec.User),
				pkgcorev1.WithName(spec.Name),
				pkgcorev1.WithCommand(spec.Cmd),
				pkgcorev1.WithResources(spec.Resources),
				pkgcorev1.WithImagePullPolicy(corev1.PullAlways),
			),
		),
	)

	_, err := typedbatchv1.Jobs(ns).Create(&jobSpec)

	return err
}

func cleanOldJobs(jobName string, ns string, history []string,
	typedbatchv1 typedbatchv1.BatchV1Interface, typedcorev1 typedcorev1.CoreV1Interface) {
	currentList := getJobList(jobName, ns, typedbatchv1)
	deleteCount := getDeleteCount(history, len(currentList))
	for i, job := range currentList {
		if i == deleteCount {
			break
		}
		if err := typedbatchv1.Jobs(ns).Delete(job.ObjectMeta.Name, &metav1.DeleteOptions{}); err != nil {
			log.Printf("Unable to delete old jobs %s", err)
		}
		if err := typedcorev1.Pods(ns).DeleteCollection(
			&metav1.DeleteOptions{},
			metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", "job-name", job.ObjectMeta.Name)},
		); err != nil {
			log.Printf("Unable to delete old pods %s", err)
		}
	}
}

func getJobList(jobName string, ns string, typedbatchv1 typedbatchv1.BatchV1Interface) []batchv1.Job {
	nsJobList, err := typedbatchv1.Jobs(ns).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Unable to retrieve old jobs %s", err)
		return nil
	}
	var jobList []batchv1.Job

	for _, job := range nsJobList.Items {
		if strings.Contains(job.ObjectMeta.Name, jobName) &&
			job.Status.Active == int32(0) {
			jobList = append(jobList, job)
		}
	}
	sort.SliceStable(jobList, jobListSortFunc(jobList))
	return jobList
}

func jobListSortFunc(jobList []batchv1.Job) func(i int, j int) bool {
	return func(i, j int) bool {
		a := jobList[i]
		b := jobList[j]

		// By timestamp
		aTime := a.ObjectMeta.CreationTimestamp
		bTime := b.ObjectMeta.CreationTimestamp
		return aTime.Before(&bTime)
	}
}

func getDeleteCount(history []string, currentCount int) int {
	var preserveCount int
	if len(history) == 0 {
		preserveCount = 3
	} else {
		preserveCount, _ = strconv.Atoi(history[0])
	}
	if currentCount > preserveCount {
		return currentCount - preserveCount
	}
	return 0
}

func getLabelsFromQP(qp map[string][]string) []types.KV {
	labels := qp["labels"]
	if len(labels) == 0 {
		return nil
	}
	var kvs []types.KV
	for _, label := range strings.Split(labels[0], ",") {
		kv := strings.Split(label, "=")
		kvs = append(kvs, types.KV{kv[0], kv[1]})
	}
	return kvs
}
