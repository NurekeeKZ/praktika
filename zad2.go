package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = make(map[int]Task)
var id = 1

// GET /tasks, POST /tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var list []Task
		for _, t := range tasks {
			list = append(list, t)
		}
		json.NewEncoder(w).Encode(list)
		return
	}

	if r.Method == "POST" {
		var t Task
		json.NewDecoder(r.Body).Decode(&t)

		t.ID = id
		id++

		tasks[t.ID] = t

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// /tasks/{id}
func taskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	switch r.Method {

	case "GET":
		json.NewEncoder(w).Encode(task)

	case "PUT":
		var updated Task
		json.NewDecoder(r.Body).Decode(&updated)

		task.Title = updated.Title
		task.Done = updated.Done
		tasks[taskID] = task

		json.NewEncoder(w).Encode(task)

	case "DELETE":
		delete(tasks, taskID)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
