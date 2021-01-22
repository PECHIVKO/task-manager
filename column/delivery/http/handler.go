package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/PECHIVKO/task-manager/column"
	"github.com/go-chi/chi"
)

type Handler struct {
	useCase column.UseCase
}

func NewHandler(useCase column.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Name      string `json:"column_name"`
	ProjectID int    `json:"project_id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	input := new(createInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.useCase.CreateColumn(r.Context(), input.Name, input.ProjectID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Column successfully created"})
	}
}

type updateNameInput struct {
	ID   int    `json:"column_id"`
	Name string `json:"column_name"`
}

func (h *Handler) UpdateName(w http.ResponseWriter, r *http.Request) {
	input := new(updateNameInput)
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
	err = h.useCase.UpdateColumnName(r.Context(), input.Name, input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Column successfully updated"})
	}
}

func (h *Handler) Move(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	pos, err := strconv.Atoi(chi.URLParam(r, "pos"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = h.useCase.MoveColumnToPosition(r.Context(), id, pos)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Column successfully updated"})
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = h.useCase.DeleteColumn(r.Context(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusNoContent, map[string]string{"message": "Column successfully deleted"})
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	col, err := h.useCase.GetColumn(r.Context(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, col)
	}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(r, "project_id"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	cols, err := h.useCase.FetchColumns(r.Context(), projectID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, cols)
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
