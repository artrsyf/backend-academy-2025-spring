package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	api "example.com/seminar02/api/oapi"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type taskManager struct {
	tasks map[string]api.Task
}

func (m *taskManager) GetTasks(w http.ResponseWriter, _ *http.Request) {
	tasks := make([]api.Task, 0, len(m.tasks))

	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}

	log.Printf("returning tasks counter [%d]\n", len(tasks))

	data, _ := json.Marshal(tasks)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (m *taskManager) PostTasks(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := api.Task{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now()
	id := uuid.New()

	task.Id = &id
	task.CreatedAt = &now
	task.UpdatedAt = &now

	m.tasks[task.Id.String()] = task

	log.Printf("task created [%s]\n", task.Id)

	data, _ := json.Marshal(&task)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (m *taskManager) GetTasksTaskId(w http.ResponseWriter, r *http.Request, taskId openapi_types.UUID) {
	task, ok := m.tasks[taskId.String()]
	if !ok {
		log.Printf("task [%s] not found\n", taskId)

		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("task was found [%s]\n", taskId)

	data, _ := json.Marshal(&task)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (m *taskManager) DeleteTasksTaskId(w http.ResponseWriter, r *http.Request, taskId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (m *taskManager) PutTasksTaskId(w http.ResponseWriter, r *http.Request, taskId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}
