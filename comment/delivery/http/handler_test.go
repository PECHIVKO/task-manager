package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PECHIVKO/task-manager/comment/usecase"
	"github.com/PECHIVKO/task-manager/models"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {

	r := chi.NewRouter()

	inp := &createInput{
		Comment: "Test Comment",
		TaskID:  1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.CommentUseCaseMock)

	uc.On("CreateComment", inp.Comment, inp.TaskID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/comments", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comments", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestHandler_Update(t *testing.T) {

	r := chi.NewRouter()

	inp := &updateInput{
		ID:      1,
		Comment: "Test Comment",
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.CommentUseCaseMock)

	uc.On("UpdateComment", inp.Comment, inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/comments", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/comments/1", bytes.NewBuffer(body))
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

	uc := new(usecase.CommentUseCaseMock)

	uc.On("DeleteComment", inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/comments", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/1", bytes.NewBuffer(body))
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

	uc := new(usecase.CommentUseCaseMock)

	com := &models.Comment{
		ID:      1,
		Task:    1,
		Date:    time.Now(),
		Comment: "Test Comment",
	}

	uc.On("GetComment", inp.ID).Return(com, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/comments", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/comments/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(com)
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

	uc := new(usecase.CommentUseCaseMock)
	coms := make([]*models.Comment, 5)
	for i := 0; i < 5; i++ {
		coms[i] = &models.Comment{
			ID:      i,
			Task:    i,
			Date:    time.Now(),
			Comment: "Test Comment",
		}
	}

	uc.On("FetchComments", inp.ID).Return(coms, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/comments", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/comments/task/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(coms)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedOutBody), w.Body.String())
}
