package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func mountRoutes(r chi.Router, deps Deps) {
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/auth", deps.Auth.Routes())
	})
}
