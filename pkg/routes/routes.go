package routes

import (
	"context"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	pkghandlers "github.com/itsmurugappan/http-handlers/pkg/handlers/favicon"

	"github.com/itsmurugappan/job-trigger/pkg/function"
)

func GetRoutes(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(300 * time.Second))

	r.Post("/", function.Context(ctx).CreateJobAndTrigger)
	r.Get("/", function.Context(ctx).CreateJobAndTrigger)
	r.Post("/{jobName}/status", function.Context(ctx).CheckJobStatus)
	r.Get("/{jobName}/status", function.Context(ctx).CheckJobStatus)
	r.Get("/favicon.ico", pkghandlers.FaviconHandler)

	return r
}
