package function

import (
	"testing"

	"gotest.tools/assert"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
)

func TestQPOverriding(t *testing.T) {
	for _, tc := range []struct {
		name  string
		want  *kubernetes.ContainerSpec
		input *kubernetes.ContainerSpec
		qp    map[string][]string
	}{{
		name:  "With overriden name",
		input: &kubernetes.ContainerSpec{Name: "foo"},
		want:  &kubernetes.ContainerSpec{Name: "bar"},
		qp:    map[string][]string{"JS_Name": []string{"bar"}},
	}, {
		name:  "With multiple overriden values",
		input: &kubernetes.ContainerSpec{Name: "foo", Image: "docker.com/muru:v1", Port: 8080, User: 1001},
		want:  &kubernetes.ContainerSpec{Name: "bar", Image: "docker.com/muru:v2", Port: 8081, User: 1002},
		qp: map[string][]string{
			"JS_Name":  []string{"bar"},
			"JS_Image": []string{"docker.com/muru:v2"},
			"JS_Port":  []string{"8081"},
			"JS_User":  []string{"1002"},
		},
	}, {
		name:  "With no overriden values",
		input: &kubernetes.ContainerSpec{Name: "foo", Image: "docker.com/muru:v1", Port: 8080, User: 1001},
		want:  &kubernetes.ContainerSpec{Name: "foo", Image: "docker.com/muru:v1", Port: 8080, User: 1001},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			err := overrideWithQP(tc.input, tc.qp, "JS_")
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.want, tc.input)
		})
	}
}
