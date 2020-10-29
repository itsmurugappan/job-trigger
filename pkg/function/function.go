package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/google/uuid"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	logging "knative.dev/pkg/logging"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/knative"
	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
	pkgbatchv1 "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes/batchv1"
	pkgcorev1 "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes/corev1"
)

func (ctx *fnCtx) CreateJobAndTrigger(w http.ResponseWriter, r *http.Request) {
	qp := map[string][]string(r.URL.Query())
	spec := specFromContext(ctx.c)
	jobSpec, err := constructJobSpec(ctx.c, qp)
	if err != nil {
		render.Respond(w, r, &Status{"", err.Error()})
		return
	}

	// delete old jobs
	cleanOldJobs(ctx.c, spec.Name, qp)

	render.Respond(w, r, createJob(ctx.c, jobSpec))
}

func (ctx *fnCtx) CheckJobStatus(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "jobName")
	ns := nsFromContext(ctx.c)

	status, err := pkgbatchv1.Client(ctx.c).GetJobStatus(*ns, name)
	if err != nil {
		render.Respond(w, r, &Status{name, err})
		return
	}
	render.Respond(w, r, &Status{name, status})
}

func createJob(ctx context.Context, jobSpec *batchv1.Job) *Status {
	ns := nsFromContext(ctx)
	job, err := pkgbatchv1.Client(ctx).CreateJob(*ns, jobSpec, true)
	if err != nil {
		return &Status{job.ObjectMeta.Name, err.Error()}
	}
	return &Status{job.ObjectMeta.Name, "Job created sucessfully!"}
}

func constructJobSpec(ctx context.Context, qp map[string][]string) (*batchv1.Job, error) {
	ns := nsFromContext(ctx)
	ksvcName := os.Getenv("K_SERVICE")
	labels := getLabelsFromQP(qp)
	originalSpec := specFromContext(ctx)
	spec := *originalSpec
	overrideWithQP(&spec, qp, "JS_")

	//get ksvc for setting owner reference
	ksvc, err := knative.Client(ctx).GetKService(*ns, ksvcName)
	if err != nil {
		logging.FromContext(ctx).Error("err getting ksvc", err)
	}

	jobSpec := pkgbatchv1.GetJob(fmt.Sprintf("%s-%s", spec.Name, uuid.New()),
		pkgbatchv1.WithTTL(int32(100)),
		pkgbatchv1.WithBackoffLimit(int32(0)),
		pkgbatchv1.WithLabels(labels),
		pkgbatchv1.WithOwnerReference(ksvc),
		pkgbatchv1.WithAnnotations([]kubernetes.KV{{"sidecar.istio.io/inject", "false"}}),
		pkgbatchv1.WithPodSpecOptions(kubernetes.PodSpec{},
			pkgcorev1.WithRestartPolicy("Never"),
			pkgcorev1.WithServiceAccount(spec.ServiceAccount),
			pkgcorev1.WithVolumes([]kubernetes.ContainerSpec{{Secrets: spec.Secrets, ConfigMaps: spec.ConfigMaps}}),
			pkgcorev1.WithContainerOptions(kubernetes.ContainerSpec{Image: spec.Image},
				pkgcorev1.WithEnv(spec.EnvVariables),
				pkgcorev1.WithEnv(pkgcorev1.GetEnvFromHTTPParam(qp)),
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
	return &jobSpec, nil
}

//overrideWithQP overrides the spec with qp values, only replaces primitive kinds
func overrideWithQP(spec interface{}, qp map[string][]string, qpPrefix string) error {
	s := reflect.ValueOf(spec)
	s = s.Elem()
	typeOfSpec := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		typ := f.Type()
		ftype := typeOfSpec.Field(i)

		overrideValue := fetchFromQP(qp, qpPrefix+ftype.Name)

		if overrideValue == nil {
			// no qp, so continue
			continue
		}

		switch typ.Kind() {
		case reflect.String:
			f.SetString(*overrideValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(*overrideValue, 0, typ.Bits())
			if err != nil {
				return err
			}

			f.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(*overrideValue, 0, typ.Bits())
			if err != nil {
				return err
			}
			f.SetUint(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(*overrideValue)
			if err != nil {
				return err
			}
			f.SetBool(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(*overrideValue, typ.Bits())
			if err != nil {
				return err
			}
			f.SetFloat(val)
		}
	}
	return nil
}

func fetchFromQP(qp map[string][]string, name string) *string {
	qpValues := qp[name]
	if len(qpValues) == 0 {
		return nil
	}
	value := qpValues[0]
	delete(qp, name)
	return &value
}

func getLabelsFromQP(qp map[string][]string) []kubernetes.KV {
	labels := qp["labels"]
	if len(labels) == 0 {
		return nil
	}
	var kvs []kubernetes.KV
	for _, label := range strings.Split(labels[0], ",") {
		kv := strings.Split(label, "=")
		kvs = append(kvs, kubernetes.KV{kv[0], kv[1]})
	}
	delete(qp, "labels")
	return kvs
}
