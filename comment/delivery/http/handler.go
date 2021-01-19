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

type updateNameInput struct {
	CommentID int    `json:"comment_id"`
	Comment   string `json:"comment"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	input := new(updateNameInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	input.CommentID, _ = strconv.Atoi(chi.URLParam(r, "id"))

	err = h.useCase.UpdateComment(r.Context(), input.Comment, input.CommentID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Comment successfully updated"})
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	commentID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.useCase.DeleteComment(r.Context(), commentID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusNoContent, map[string]string{"message": "Comment successfully deleted"})
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	commentID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	com, err := h.useCase.GetComment(r.Context(), commentID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, com)
	}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	taskID, _ := strconv.Atoi(chi.URLParam(r, "task_id"))

	coms, err := h.useCase.FetchComments(r.Context(), taskID)
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
