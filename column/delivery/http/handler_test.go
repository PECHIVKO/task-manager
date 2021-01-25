package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PECHIVKO/task-manager/column/usecase"
	"github.com/PECHIVKO/task-manager/models"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {

	r := chi.NewRouter()

	inp := &createInput{
		Name:      "Test Column",
		ProjectID: 1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ColumnUseCaseMock)

	uc.On("CreateColumn", inp.Name, inp.ProjectID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/columns", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/columns", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestHandler_UpdateName(t *testing.T) {

	r := chi.NewRouter()

	inp := &updateNameInput{
		Name: "Test Column",
		ID:   1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ColumnUseCaseMock)

	uc.On("UpdateColumnName", inp.Name, inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/columns", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/columns/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestHandler_Move(t *testing.T) {

	r := chi.NewRouter()

	inp := &moveInput{
		Position: 4,
		ID:       1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ColumnUseCaseMock)

	uc.On("MoveColumnToPosition", inp.ID, inp.Position).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/columns", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/columns/1/move/4", bytes.NewBuffer(body))
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

	uc := new(usecase.ColumnUseCaseMock)

	uc.On("DeleteColumn", inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/columns", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/columns/1", bytes.NewBuffer(body))
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

	uc := new(usecase.ColumnUseCaseMock)

	col := &models.Column{
		ID:       1,
		Name:     "Test Column",
		Project:  1,
		Position: 1,
	}

	uc.On("GetColumn", inp.ID).Return(col, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/columns", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/columns/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(col)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedOutBody), w.Body.String())
}

func TestHandler_Fetch(t *testing.T) {

	r := chi.NewRouter()

	inp := &fetchInput{
		ID: 1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.ColumnUseCaseMock)
	cols := make([]*models.Column, 5)
	for i := 0; i < 5; i++ {
		cols[i] = &models.Column{
			ID:       i,
			Name:     "Test Column",
			Project:  i,
			Position: i,
		}
	}

	uc.On("FetchColumns", inp.ID).Return(cols, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/columns", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/columns/project/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(cols)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedOutBody), w.Body.String())
}
