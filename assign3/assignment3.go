package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

	for {
		displayMenu()
		choiceStr := getUserInput()
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid choice. Please enter a number between 1 and 4.")
			continue
		}

		switch choice {
		case 1:
			fmt.Print("Enter task description: ")
			description := getUserInput()
			if description == "" {
				fmt.Println("Task description cannot be empty.")
				continue
			}
			addedTask := tracker.AddTask(description)
			fmt.Printf("Task Added: %d - %s\n", addedTask.ID, addedTask.Description)
		case 2:
			fmt.Println(tracker.ListTasks())
		case 3:
			fmt.Print("Enter ID of task to mark as completed: ")
			idStr := getUserInput()
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID. Please enter a valid number.")
				continue
			}
			_, msg := tracker.CompleteTask(id)
			fmt.Println(msg)
		case 4:
			fmt.Println("Exiting Task Tracker. Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please choose a number between 1 and 4.")
		}
	}
}
