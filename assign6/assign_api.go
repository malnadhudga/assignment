package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func httpmark(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {

	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("give the ID"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id > len(tracker.tasks) || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please Enter the valid ID"))
		return
	}

	_, message := tracker.CompleteTask(id)
	fmt.Sprintf(message)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
	return

}

func httppostTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {

	queryValues := r.URL.Query()
	task := queryValues.Get("task")

	if task == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing id parameter"))
		return
	}

	tracker.AddTask(task)
	w.WriteHeader(http.StatusCreated) //Added the status code
	return
}

func httpListtask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	tasks := tracker.ListTasks()

	fmt.Sprintf(tasks)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tasks))
	return
}

func httpDelete(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {

	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")

	if idStr == "" {
		//w.Write([]byte("give the ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id > len(tracker.tasks) || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please Enter the valid ID"))

		return
	}

	id = id - 1
	t1 := tracker.tasks[:id]
	t2 := tracker.tasks[id+1:]
	tracker.tasks = append(t1, t2...)

	tasks := tracker.ListTasks()

	fmt.Sprintf(tasks)
	w.WriteHeader(http.StatusNoContent) //Added the status code
	w.Write([]byte(tasks))
	return
}

func httpListbyId(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if id > len(tracker.tasks) || id < 0 {
		w.Write([]byte("Please Enter the valid ID"))
		return
	}

	flag := false

	for _, task := range tracker.tasks {
		if id == task.ID {
			flag = true
			fmt.Sprintf(task.Description)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(task.Description))
			return
		}
	}

	if !flag {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Task not found ofr this ID"))
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
// It uses a pointer receiver (*TaskTracker) because it modifies the state of a Task within the tracker's slice.
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

// displayMenu prints the interactive menu options to the console.
func displayMenu() {
	fmt.Println("\n--- Personal Task Tracker ---")
	fmt.Println("1. Add a new task")
	fmt.Println("2. List all pending tasks")
	fmt.Println("3. Mark a task as completed")
	fmt.Println("4. Exit")
	fmt.Print("Choose an option: ")
}

// getUserInput reads a line of text from the standard input.
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// main function orchestrates the CLI interaction.
func main() {
	tracker := NewTaskTracker()

	http.HandleFunc("GET /task", func(w http.ResponseWriter, r *http.Request) { httpListtask(w, r, tracker) })
	http.HandleFunc("GET /task/{id}", func(w http.ResponseWriter, r *http.Request) { httpListbyId(w, r, tracker) })
	http.HandleFunc("POST /task", func(w http.ResponseWriter, r *http.Request) { httppostTask(w, r, tracker) })
	http.HandleFunc("PUT /task", func(w http.ResponseWriter, r *http.Request) { httpmark(w, r, tracker) })
	http.HandleFunc("DELETE /task", func(w http.ResponseWriter, r *http.Request) { httpDelete(w, r, tracker) })

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
