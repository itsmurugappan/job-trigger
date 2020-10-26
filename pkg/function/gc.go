package function

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
	logging "knative.dev/pkg/logging"
)

func cleanOldJobs(ctx context.Context, jobName string, qp map[string][]string) {
	ns := nsFromContext(ctx)
	candidates := getJobList(ctx, jobName, *ns)

	history := qp["history"]
	deleteCount := getDeleteCount(history, len(candidates))
	//
	gc(ctx, *ns, deleteCount, candidates)
	//clear query params
	delete(qp, "history")
}

func gc(ctx context.Context, ns string, deleteCount int, candidates []batchv1.Job) {
	cs := kubernetes.KubernetesCSFromContext(ctx)
	logger := logging.FromContext(ctx)
	logger.Infof("deleting last %d jobs", deleteCount)
	for i, job := range candidates {
		if i == deleteCount {
			break
		}
		if err := cs.BatchV1().Jobs(ns).Delete(ctx, job.ObjectMeta.Name, metav1.DeleteOptions{}); err != nil {
			logger.Errorf("Unable to delete old jobs %s", err)
		}
		if err := cs.CoreV1().Pods(ns).DeleteCollection(
			ctx,
			metav1.DeleteOptions{},
			metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", "job-name", job.ObjectMeta.Name)},
		); err != nil {
			logger.Errorf("Unable to delete old pods %s", err)
		}
	}
}

func getJobList(ctx context.Context, jobName, ns string) []batchv1.Job {
	cs := kubernetes.KubernetesCSFromContext(ctx)

	nsJobList, err := cs.BatchV1().Jobs(ns).List(ctx, metav1.ListOptions{})

	logger := logging.FromContext(ctx)

	if err != nil {
		logger.Errorf("Unable to retrieve old jobs %s", err)
		return nil
	}

	var jobList []batchv1.Job

	for _, job := range nsJobList.Items {
		if strings.Contains(job.ObjectMeta.Name, jobName) &&
			isGCEligible(ctx, job) {
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

func isGCEligible(ctx context.Context, job batchv1.Job) bool {
	if job.Status.Active == int32(0) {
		return true
	}
	// if it is active see if its created atleast 10 hours ago
	difference := time.Now().Sub(job.Status.StartTime.Time)
	logging.FromContext(ctx).Infof("Job %s time created %v", job.ObjectMeta.Name, difference)
	if difference > time.Duration(10*time.Hour) {
		return true
	}
	return false
}
