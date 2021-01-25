package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/PECHIVKO/task-manager/task"
	"github.com/go-chi/chi"
)

type Handler struct {
	useCase task.UseCase
}

func NewHandler(useCase task.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Name        string `json:"task_name"`
	Description string `json:"task_description"`
	ColumnID    int    `json:"column_id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	input := new(createInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.useCase.CreateTask(r.Context(), input.Name, input.Description, input.ColumnID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Task successfully created"})
	}
}

type updateInput struct {
	ID          int    `json:"task_id"`
	Name        string `json:"task_name"`
	Description string `json:"task_description"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	input := new(updateInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	input.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = h.useCase.UpdateTask(r.Context(), input.Name, input.Description, input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Task successfully updated"})
	}
}

type moveInput struct {
	ID     int `json:"task_id"`
	Column int `json:"column_id"`
}

func (h *Handler) Move(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(moveInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	input.Column, err = strconv.Atoi(chi.URLParam(r, "column"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = h.useCase.MoveToColumn(r.Context(), input.ID, input.Column)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Task successfully moved"})
	}
}

type priorityInput struct {
	ID       int `json:"task_id"`
	Priority int `json:"priority"`
}

func (h *Handler) ChangePriority(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(priorityInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	input.Priority, err = strconv.Atoi(chi.URLParam(r, "priority"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = h.useCase.ChangeTaskPriority(r.Context(), input.ID, input.Priority)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Task successfully updated"})
	}

}

type deleteInput struct {
	ID int `json:"task"`
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(deleteInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = h.useCase.DeleteTask(r.Context(), input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusNoContent, map[string]string{"message": "Task successfully deleted"})
	}

}

type getInput struct {
	ID int `json:"task_id"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(getInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	task, err := h.useCase.GetTask(r.Context(), input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, task)
	}

}

type fetchInput struct {
	ID int `json:"column_id"`
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(fetchInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "column_id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	tasks, err := h.useCase.FetchTasks(r.Context(), input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, tasks)
	}

}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"Error": msg})
}
