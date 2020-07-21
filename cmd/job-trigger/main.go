package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	errs "github.com/itsmurugappan/http-handlers/pkg/errors"
	"github.com/itsmurugappan/http-handlers/pkg/handlers/favicon"
	"github.com/itsmurugappan/job-trigger/pkg/function"
	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
)

func main() {
	http.HandleFunc("/", jobHandler)
	http.HandleFunc("/favicon.ico", favicon.FaviconHandler)
	http.ListenAndServe(":8080", nil)
}

func jobHandler(w http.ResponseWriter, r *http.Request) {
	cs, err := kubernetes.GetInClusterKubeClient()
	if err != nil {
		errs.HandleErrors(w, err)
		return
	}

	raw, present := os.LookupEnv("spec")
	if !present {
		errs.HandleErrors(w, fmt.Errorf("'spec' not provided as a env variable"))
		return
	}

	containerSpec := kubernetes.ContainerSpec{}
	if err := json.Unmarshal([]byte(raw), &containerSpec); err != nil {
		errs.HandleErrors(w, err)
		return
	}

	if err := function.CreateJob(containerSpec, kubernetes.GetCurrentNamespace(), map[string][]string(r.URL.Query()), cs.CoreV1(), cs.BatchV1()); err != nil {
		errs.HandleErrors(w, err)
		return
	}

	fmt.Fprintf(w, "Job Created\n")
}
