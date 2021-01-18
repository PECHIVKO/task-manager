package http

import (
	"github.com/PECHIVKO/task-manager/task"
	"github.com/go-chi/chi"
)

func Routes(uc task.UseCase) *chi.Mux {
	router := chi.NewRouter()
	h := NewHandler(uc)

	router.Get("/{id:[0-9]+}", h.Get)
	router.Post("/", h.Create)
	router.Delete("/{id:[0-9]+}", h.Delete)
	router.Get("/column/{column_id:[0-9]+}", h.Fetch)
	router.Put("/{id:[0-9]+}", h.Update)
	router.Put("/{id:[0-9]+}/priority/{priority:[0-9]+}", h.ChangePriority)
	router.Put("/{id:[0-9]+}/move/{column:[0-9]+}", h.Move)
	return router
}
