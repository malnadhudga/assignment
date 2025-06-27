package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func httppostTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	queryValues := r.URL.Query()
	task := queryValues.Get("task")
	if task == "" {
		w.Write([]byte("Missing id parameter"))
		return
	}
	tracker.AddTask(task)

}
func httpListallTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	tasks := tracker.ListTasks()
	fmt.Sprintf(tasks)
	w.Write([]byte(tasks))

}
func httpDeletetask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	id = id - 1
	if id < 0 || id >= len(tracker.ListTasks()) {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tracker.tasks = append(tracker.tasks[:id], tracker.tasks[id+1:]...)
	tasks := tracker.ListTasks()
	fmt.Sprintf(tasks)
	w.Write([]byte(tasks))

}
func httpListbyId(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	flag := false
	id = id - 1
	for i, task := range tracker.tasks {
		if id == i {
			flag = true
			fmt.Sprintf(task.Description)
			w.Write([]byte(task.Description))
		}
	}
	if !flag {
		w.Write([]byte("Task not found ofr this ID"))
	}

}
func httpPutTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	_, msg := tracker.CompleteTask(id)
	w.Write([]byte(msg))
}

type Task struct {
	ID          int
	Description string
	Completed   bool
}

type TaskTracker struct {
	tasks     []Task
	nextIDGen func() int
}

func idGenerator() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}
func NewTaskTracker() *TaskTracker {
	return &TaskTracker{
		tasks:     []Task{},
		nextIDGen: idGenerator(),
	}
}
func (tt *TaskTracker) AddTask(description string) Task {
	newID := tt.nextIDGen()
	newTask := Task{
		ID:          newID,
		Description: description,
		Completed:   false,
	}
	tt.tasks = append(tt.tasks, newTask)
	return newTask
}

func (tt *TaskTracker) ListTasks() string {
	s := "Pending Tasks:\n"
	foundPending := false
	for _, task := range tt.tasks {
		if !task.Completed {
			s += fmt.Sprintf("%d: %s\n", task.ID, task.Description)
			foundPending = true
		}
	}
	if !foundPending {
		s += "No pending tasks."
	}
	return s
}

func (tt *TaskTracker) CompleteTask(id int) (bool, string) {
	for i := range tt.tasks {
		if tt.tasks[i].ID == id {
			if tt.tasks[i].Completed {
				return false, fmt.Sprintf("Task %d is already completed.", id)
			}
			tt.tasks[i].Completed = true
			return true, fmt.Sprintf("Marking task %d as completed: %s", id, tt.tasks[i].Description)
		}
	}
	return false, fmt.Sprintf("Task with ID %d not found.", id)
}

// main function orchestrates the CLI interaction.
func main() {
	tracker := NewTaskTracker()
	http.HandleFunc("GET /task", func(w http.ResponseWriter, r *http.Request) { httpListallTask(w, r, tracker) })
	http.HandleFunc("GET /task/{id}", func(w http.ResponseWriter, r *http.Request) { httpListbyId(w, r, tracker) })
	http.HandleFunc("POST /task", func(w http.ResponseWriter, r *http.Request) { httppostTask(w, r, tracker) })
	http.HandleFunc("DELETE /task", func(w http.ResponseWriter, r *http.Request) { httpDeletetask(w, r, tracker) })
	http.HandleFunc("PATCH /task", func(w http.ResponseWriter, r *http.Request) { httpPutTask(w, r, tracker) })
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
