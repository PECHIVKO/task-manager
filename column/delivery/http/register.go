package http

import (
	"github.com/PECHIVKO/task-manager/column"
	"github.com/go-chi/chi"
)

func Routes(uc column.UseCase) *chi.Mux {
	router := chi.NewRouter()
	h := NewHandler(uc)

	router.Get("/{id:[0-9]+}", h.Get)
	router.Post("/create", h.Create)
	router.Delete("/{id:[0-9]+}", h.Delete)
	router.Get("/project/{project_id:[0-9]+}", h.Fetch)
	router.Put("/{id:[0-9]+}", h.UpdateName)
	router.Put("/move/{id:[0-9]+}/{pos:[0-9]+}", h.Move)
	return router
}
