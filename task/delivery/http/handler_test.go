package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/PECHIVKO/task-manager/task/usecase"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {

	r := chi.NewRouter()

	inp := &createInput{
		Name:        "Test Task",
		Description: "Test Description",
		ColumnID:    1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.TaskUseCaseMock)

	uc.On("CreateTask", inp.Name, inp.Description, inp.ColumnID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/tasks", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestHandler_Update(t *testing.T) {

	r := chi.NewRouter()

	inp := &updateInput{
		Name:        "Test Task",
		Description: "Test Description",
		ID:          1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.TaskUseCaseMock)

	uc.On("UpdateTask", inp.Name, inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/tasks", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestHandler_Move(t *testing.T) {

	r := chi.NewRouter()

	inp := &moveInput{
		Column: 4,
		ID:     1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.TaskUseCaseMock)

	uc.On("MoveToColumn", inp.ID, inp.Column).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/tasks", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/1/move/4", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestHandler_ChangePriority(t *testing.T) {

	r := chi.NewRouter()

	inp := &priorityInput{
		Priority: 4,
		ID:       1,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc := new(usecase.TaskUseCaseMock)

	uc.On("ChangeTaskPriority", inp.ID, inp.Priority).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/tasks", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/1/priority/4", bytes.NewBuffer(body))
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

	uc := new(usecase.TaskUseCaseMock)

	uc.On("DeleteTask", inp.ID).Return(nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/tasks", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", bytes.NewBuffer(body))
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

	uc := new(usecase.TaskUseCaseMock)

	task := &models.Task{
		ID:          1,
		Column:      1,
		Priority:    1,
		Name:        "Test Task",
		Description: "Test Description",
	}

	uc.On("GetTask", inp.ID).Return(task, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/tasks", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(task)
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

	uc := new(usecase.TaskUseCaseMock)
	tasks := make([]*models.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = &models.Task{
			ID:          i,
			Column:      i,
			Priority:    i,
			Name:        "Test Task",
			Description: "Test Description",
		}
	}

	uc.On("FetchTasks", inp.ID).Return(tasks, nil)

	r.Route("/", func(r chi.Router) {
		r.Mount("/tasks", Routes(uc))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/column/1", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	expectedOutBody, err := json.Marshal(tasks)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedOutBody), w.Body.String())
}
