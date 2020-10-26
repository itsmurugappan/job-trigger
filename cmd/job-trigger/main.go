package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	logging "knative.dev/pkg/logging"

	"github.com/itsmurugappan/job-trigger/pkg/function"
	"github.com/itsmurugappan/job-trigger/pkg/routes"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx := function.BuildContext()

	logger := logging.FromContext(ctx)
	defer logger.Sync()

	if err := function.ValidateK8sResource(ctx); err != nil {
		panic(err)
	}

	r := routes.GetRoutes(ctx)

	go func() {
		http.ListenAndServe(":8080", r)
	}()

	<-done
	logger.Info("Shutting down..")
}
