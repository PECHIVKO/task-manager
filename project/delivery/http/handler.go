package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/PECHIVKO/task-manager/project"
	"github.com/go-chi/chi"
)

type Handler struct {
	useCase project.UseCase
}

func NewHandler(useCase project.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type createInput struct {
	Name        string `json:"project_name"`
	Description string `json:"project_description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	input := new(createInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.useCase.CreateProject(r.Context(), input.Name, input.Description)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Project successfully created"})
	}
}

type updateInput struct {
	ID          int    `json:"project_id"`
	Name        string `json:"project_name"`
	Description string `json:"project_description"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	input := new(updateInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	input.ID, _ = strconv.Atoi(chi.URLParam(r, "id"))

	err = h.useCase.UpdateProject(r.Context(), input.Name, input.Description, input.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "Project successfully updated"})
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	err := h.useCase.DeleteProject(r.Context(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusNoContent, map[string]string{"message": "Project successfully deleted"})
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	pr, err := h.useCase.GetProject(r.Context(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, pr)
	}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	prs, err := h.useCase.FetchProjects(r.Context())
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondwithJSON(w, http.StatusOK, prs)
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
