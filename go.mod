module github.com/itsmurugappan/job-trigger

go 1.15

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.1.1
	github.com/itsmurugappan/http-handlers v0.0.0-20200524185756-1c2f336610f6
	github.com/itsmurugappan/kubernetes-resource-builder v0.0.0-20201026014106-cb22ee8529e1
	go.uber.org/zap v1.16.0 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	knative.dev/pkg v0.0.0-20200922164940-4bf40ad82aab
)

replace (
	k8s.io/api => k8s.io/api v0.18.8

	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
)
