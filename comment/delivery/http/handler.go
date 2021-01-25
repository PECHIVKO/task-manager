package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/PECHIVKO/task-manager/comment"
	"github.com/go-chi/chi"
)

type Handler struct {
	useCase comment.UseCase
}

func NewHandler(useCase comment.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Comment string `json:"comment"`
	TaskID  int    `json:"task_id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	input := new(createInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.useCase.CreateComment(r.Context(), input.Comment, input.TaskID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Comment successfully created"})
	}
}

type updateInput struct {
	ID      int    `json:"comment_id"`
	Comment string `json:"comment"`
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
	err = h.useCase.UpdateComment(r.Context(), input.Comment, input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Comment successfully updated"})
	}
}

type deleteInput struct {
	ID int `json:"comment_id"`
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(deleteInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = h.useCase.DeleteComment(r.Context(), input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusNoContent, map[string]string{"message": "Comment successfully deleted"})
	}
}

type getInput struct {
	ID int `json:"comment_id"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(getInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	com, err := h.useCase.GetComment(r.Context(), input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, com)
	}
}

type fetchInput struct {
	ID int `json:"task_id"`
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	var err error
	input := new(fetchInput)
	input.ID, err = strconv.Atoi(chi.URLParam(r, "task_id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	coms, err := h.useCase.FetchComments(r.Context(), input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, coms)
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
