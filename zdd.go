package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// =========================
// 1. SUM / MIN / MAX
// =========================
func calcStats(nums []int) (sum, min, max int, err error) {
	if len(nums) == 0 {
		return 0, 0, 0, fmt.Errorf("empty slice")
	}

	sum = 0
	min = nums[0]
	max = nums[0]

	for _, n := range nums {
		sum += n

		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	return sum, min, max, nil
}

// =========================
// 2. CALCULATOR
// =========================
func calculator(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}

// =========================
// 3. WORKER (GOROUTINES)
// =========================
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("Worker", id, "processing job", j)
		results <- j * j
	}
}

// =========================
// TASK API
// =========================
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

	switch r.Method {

	case "GET":
		var list []Task
		for _, t := range tasks {
			list = append(list, t)
		}
		json.NewEncoder(w).Encode(list)

	case "POST":
		var t Task
		json.NewDecoder(r.Body).Decode(&t)

		t.ID = id
		id++

		tasks[t.ID] = t

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

// =========================
// MAIN
// =========================
func main() {

	// -------- TASK 1 --------
	nums := []int{5, 2, 9, 1, 7}

	sum, min, max, err := calcStats(nums)
	if err != nil {
		fmt.Println("Stats error:", err)
	} else {
		fmt.Println("=== STATS ===")
		fmt.Println("Sum:", sum)
		fmt.Println("Min:", min)
		fmt.Println("Max:", max)
	}

	// -------- TASK 2 --------
	fmt.Println("\n=== CALCULATOR ===")

	if res, err := calculator(10, 5, "+"); err == nil {
		fmt.Println("10 + 5 =", res)
	}

	if res, err := calculator(10, 0, "/"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(res)
	}

	// -------- TASK 3 --------
	fmt.Println("\n=== WORKERS ===")

	jobs := make(chan int)
	results := make(chan int)

	for i := 1; i <= 3; i++ {
		go worker(i, jobs, results)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	for i := 1; i <= 5; i++ {
		fmt.Println("Result:", <-results)
	}

	// -------- HTTP SERVER --------
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", taskHandler)

	fmt.Println("\nServer started on :8080")
	http.ListenAndServe(":8080", nil)
}
