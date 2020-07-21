package function

import (
	"fmt"

	"testing"

	"gotest.tools/assert"

	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"

	types "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
	constants "github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes/corev1"
	teststubcorev1 "github.com/itsmurugappan/kubernetes-resource-builder/pkg/test/kubernetes/corev1"
)

func TestCreateJob(t *testing.T) {
	for _, tc := range []struct {
		name           string
		want           string
		inputModel     types.ContainerSpec
		ns             string
		queryParams    map[string][]string
		runtimeObjects []runtime.Object
	}{{
		name:       "Create fresh job",
		want:       "",
		inputModel: types.ContainerSpec{Name: "foo"},
		ns:         "foo",
	}, {
		name:       "Create job without cm",
		want:       fmt.Sprintf(constants.CM_MISSING, "cm1", "foo"),
		inputModel: types.ContainerSpec{Name: "foo", ConfigMaps: teststubcorev1.ConstructMounts([]string{"cm1"}, []string{"p1"})},
		ns:         "foo",
	}, {
		name:       "Create job without secret",
		want:       fmt.Sprintf(constants.SECRET_MISSING, "s1", "foo"),
		inputModel: types.ContainerSpec{Name: "foo", Secrets: teststubcorev1.ConstructMounts([]string{"s1"}, []string{"p1"})},
		ns:         "foo",
	}, {
		name:       "Create job without secret in env from",
		want:       fmt.Sprintf(constants.SECRET_MISSING, "s1", "foo"),
		inputModel: types.ContainerSpec{Name: "foo", EnvFromSecretorCM: []types.EnvFrom{{"s1", "Secret"}}},
		ns:         "foo",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			err := CreateJob(tc.inputModel,
				tc.ns,
				tc.queryParams,
				testclient.NewSimpleClientset((tc.runtimeObjects)...).CoreV1(),
				testclient.NewSimpleClientset().BatchV1(),
			)
			if tc.want == "" {
				assert.NilError(t, err)
			} else {
				assert.Error(t, err, tc.want)
			}
		})
	}
}
