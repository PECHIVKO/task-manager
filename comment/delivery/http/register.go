package http

import (
	"github.com/PECHIVKO/task-manager/comment"
	"github.com/go-chi/chi"
)

func Routes(uc comment.UseCase) *chi.Mux {
	router := chi.NewRouter()
	h := NewHandler(uc)

	router.Get("/{id:[0-9]+}", h.Get)
	router.Post("/", h.Create)
	router.Delete("/{id:[0-9]+}", h.Delete)
	router.Get("/task/{task_id:[0-9]+}", h.Fetch)
	router.Put("/{id:[0-9]+}", h.Update)
	return router
}
