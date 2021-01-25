package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/PECHIVKO/task-manager/project/usecase"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {

	r := chi.NewRouter()

	inp := &createInput{
		Name:        "Test Project",
		Description: "Test Description",
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ProjectUseCaseMock)

	uc.On("CreateProject", inp.Name, inp.Description).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/projects", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestHandler_UpdateName(t *testing.T) {

	r := chi.NewRouter()

	inp := &updateInput{
		Description: "Test Description",
		Name:        "Test Project",
		ID:          1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ProjectUseCaseMock)

	uc.On("UpdateProject", inp.Name, inp.Description, inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/projects", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/projects/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestHandler_Delete(t *testing.T) {

	r := chi.NewRouter()

	inp := &deleteInput{
		ID: 1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ProjectUseCaseMock)

	uc.On("DeleteProject", inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/projects", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/projects/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestHandler_Get(t *testing.T) {

	r := chi.NewRouter()

	inp := &getInput{
		ID: 1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ProjectUseCaseMock)

	pr := &models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test Description",
	}

	uc.On("GetProject", inp.ID).Return(pr, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/projects", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(pr)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedOutBody), w.Body.String())
}

func TestHandler_Fetch(t *testing.T) {

	r := chi.NewRouter()

	uc := new(usecase.ProjectUseCaseMock)
	cols := make([]*models.Project, 5)
	for i := 0; i < 5; i++ {
		cols[i] = &models.Project{
			ID:          i,
			Name:        "Test Project",
			Description: "Test Description",
		}
	}

	uc.On("FetchProjects").Return(cols, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/projects", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects", nil)
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(cols)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedOutBody), w.Body.String())
}
