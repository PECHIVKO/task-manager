package http

import (
	"github.com/PECHIVKO/task-manager/project"
	"github.com/go-chi/chi"
)

func Routes(uc project.UseCase) *chi.Mux {
	router := chi.NewRouter()
	h := NewHandler(uc)

	router.Get("/{id:[0-9]+}", h.Get)
	router.Post("/", h.Create)
	router.Delete("/{id:[0-9]+}", h.Delete)
	router.Get("/", h.Fetch)
	router.Put("/{id:[0-9]+}", h.Update)
	return router
}
