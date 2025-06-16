package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Task represents a single unit of work with a unique identifier,
// a descriptive string, and a completion status...
type Task struct {
	ID          int
	Description string
	Completed   bool
}

// TaskTracker manages a collection of tasks, providing functionalities
// to add, list, and mark tasks as complete. It also handles the generation
// of unique task IDs.
type TaskTracker struct {
	tasks     []Task
	nextIDGen func() int
}

// idGenerator returns a closure that produces unique, sequential integer IDs
// starting from 1.
//
// Returns:
//
//	func() int: A function that, when called, returns the next unique integer ID.
func idGenerator() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

// NewTaskTracker creates and initializes a new TaskTracker instance.
//
// Returns:
//
//	*TaskTracker: A pointer to a newly created and initialized TaskTracker.
func NewTaskTracker() *TaskTracker {
	return &TaskTracker{
		tasks:     []Task{},
		nextIDGen: idGenerator(),
	}
}

// AddTask appends a new task with the given description to the tracker.
// The task is assigned a unique ID and initialized as incomplete.
//
// Parameters:
//
//	description (string): The textual description of the task to be added.
func (tt *TaskTracker) AddTask(description string) {
	newID := tt.nextIDGen()
	newTask := Task{
		ID:          newID,
		Description: description,
		Completed:   false,
	}
	tt.tasks = append(tt.tasks, newTask)
	fmt.Printf("Task Added: %d - %s\n", newTask.ID, newTask.Description)
}

// ListTasks iterates through all tasks managed by the tracker and
// prints only those that are currently incomplete to standard output.
// If no pending tasks are found, a corresponding message is displayed.
func (tt *TaskTracker) ListTasks() {
	fmt.Println("\nPending Tasks:")
	foundPending := false
	for _, task := range tt.tasks {
		if !task.Completed {
			fmt.Printf("%d: %s\n", task.ID, task.Description)
			foundPending = true
		}
	}
	if !foundPending {
		fmt.Println("No pending tasks.")
	}
}

// CompleteTask marks the task with the specified ID as completed.
// If the task is not found or is already completed, appropriate messages
// are displayed to standard output.
//
// Parameters:
//
//	id (int): The unique identifier of the task to be marked as completed.
func (tt *TaskTracker) CompleteTask(id int) {
	taskFound := false
	for i := range tt.tasks {
		if tt.tasks[i].ID == id {
			if tt.tasks[i].Completed {
				fmt.Printf("Task %d is already completed.\n", id)
			} else {
				tt.tasks[i].Completed = true
				fmt.Printf("Marking task %d as completed: %s\n", id, tt.tasks[i].Description)
			}
			taskFound = true
			break
		}
	}
	if !taskFound {
		fmt.Printf("Task with ID %d not found.\n", id)
	}
}

// displayMenu presents the user with the available options for interacting
// with the task tracker application by printing them to standard output.
func displayMenu() {
	fmt.Println("\n--- Personal Task Tracker ---")
	fmt.Println("1. Add a new task")
	fmt.Println("2. List all pending tasks")
	fmt.Println("3. Mark a task as completed")
	fmt.Println("4. Exit")
	fmt.Print("Choose an option: ")
}

// getUserInput reads a single line of text from standard input,
// trims leading/trailing whitespace, and returns the resulting string.
//
// Returns:
//
//	string: The trimmed string input from the user.
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// main is the entry point of the Task Tracker application.
// It initializes a new TaskTracker and enters a loop to display the menu,
// process user input, and perform the requested task operations until
// the user chooses to exit.
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
			tracker.AddTask(description)
		case 2:
			tracker.ListTasks()
		case 3:
			fmt.Print("Enter ID of task to mark as completed: ")
			idStr := getUserInput()
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID. Please enter a valid number.")
				continue
			}
			tracker.CompleteTask(id)
		case 4:
			fmt.Println("Exiting Task Tracker. Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please choose a number between 1 and 4.")
		}
	}
}
