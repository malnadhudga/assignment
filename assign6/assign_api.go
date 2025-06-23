package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	errID  = errors.New("invalid ID parameter")
	errmax = errors.New("please enter a valid ID within range")
	errint = errors.New("invalid ID format")
)

func httpmark(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("give the ID"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id > len(tracker.tasks) || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		_, err = w.Write([]byte("Please Enter the valid ID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, message := tracker.CompleteTask(id)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(message))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func httppostTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	queryValues := r.URL.Query()
	task := queryValues.Get("task")

	if task == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Missing id parameter"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	tracker.AddTask(task)
	w.WriteHeader(http.StatusCreated)
}

func httpListtask(w http.ResponseWriter, _ *http.Request, tracker *TaskTracker) {
	tasks := tracker.ListTasks()

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(tasks))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseAndValidateID(r *http.Request, maxID int) (int, error) {
	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")

	if idStr == "" {
		return 0, errID
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errint
	}

	if id <= 0 || id > maxID {
		return 0, errmax
	}

	return id, nil
}

func httpDelete(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	maxID := len(tracker.tasks)

	id, err := parseAndValidateID(r, maxID)
	if err != nil {
		return
	}

	if id > len(tracker.tasks) || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Please Enter a valid ID within range"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	indexToDelete := id - 1
	if indexToDelete <= 0 && indexToDelete > len(tracker.tasks) {
		http.Error(w, "Task not found for deletion", http.StatusNotFound) // Use 404 for not found
		return
	}

	if len(tracker.tasks) == 1 {
		tracker.tasks = []Task{}
	} else {
		tracker.tasks = append(tracker.tasks[:indexToDelete], tracker.tasks[indexToDelete+1:]...)
	}

	tasks := tracker.ListTasks()

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Task deleted successfully. Updated list:\n" + tasks))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func httpListbyID(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if id > len(tracker.tasks) || id < 0 {
		_, err = w.Write([]byte("Please Enter the valid ID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	for _, task := range tracker.tasks {
		if id == task.ID {
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(task.Description))

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write([]byte("Task not found ofr this ID"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Task represents a single task in our tracker.
type Task struct {
	ID          int
	Description string
	Completed   bool
}

// TaskTracker manages the collection of tasks and generates unique IDs.
type TaskTracker struct {
	tasks     []Task
	nextIDGen func() int
}

// idGenerator is a closure that generates unique sequential integer IDs.
// It encapsulates the 'id' counter, so it's not a global variable.
func idGenerator() func() int {
	id := 0

	return func() int {
		id++
		return id
	}
}

// NewTaskTracker creates and initializes a new TaskTracker instance.
// It also sets up the unique ID generator.
func NewTaskTracker() *TaskTracker {
	return &TaskTracker{
		tasks:     []Task{},
		nextIDGen: idGenerator(),
	}
}

// AddTask adds a new task to the tracker and returns the added Task.
// It uses a pointer receiver (*TaskTracker) because it modifies the TaskTracker's state (its 'tasks' slice).
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

// ListTasks displays all pending tasks.
// It uses a pointer receiver (*TaskTracker) because it operates on the TaskTracker's 'tasks' slice,
// even though it doesn't modify it directly in this function (good practice for methods operating on collections).
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

// CompleteTask marks a task as completed given its ID.
// It returns a boolean indicating if the task was found and its completion status was changed,
// and a string message describing the outcome.
func (tt *TaskTracker) CompleteTask(id int) (success bool, message string) {
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

func main() {
	tracker := NewTaskTracker()

	http.HandleFunc("GET /task", func(w http.ResponseWriter, r *http.Request) { httpListtask(w, r, tracker) })
	http.HandleFunc("GET /task/{id}", func(w http.ResponseWriter, r *http.Request) { httpListbyID(w, r, tracker) })
	http.HandleFunc("POST /task", func(w http.ResponseWriter, r *http.Request) { httppostTask(w, r, tracker) })
	http.HandleFunc("PUT /task", func(w http.ResponseWriter, r *http.Request) { httpmark(w, r, tracker) })
	http.HandleFunc("DELETE /task", func(w http.ResponseWriter, r *http.Request) { httpDelete(w, r, tracker) })

	server := &http.Server{
		Addr: ":8080",

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on port %s", server.Addr)
	err := server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed to start: %v", err)
	}
}
