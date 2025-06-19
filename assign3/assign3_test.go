package main

import (
	"fmt"
	"testing"
)

func TestCompleteTask(t *testing.T) {
	tracker := NewTaskTracker()
	tracker.AddTask("Read a book")
	tracker.AddTask("Play")
	tracker.AddTask("Clean")

	success, msg := tracker.CompleteTask(2)
	if !success {
		t.Errorf("Expected task 2 to be completed successfully, but it failed: %s", msg)
	}
	fmt.Println(msg)
	if msg != "Marking task 2 as completed: Play" {
		t.Errorf("Incorrect success message. Got: '%s'", msg)
	}
	if !tracker.tasks[1].Completed {
		t.Errorf("Task 2 should be marked as completed")
	}

	success, msg = tracker.CompleteTask(2)
	if success {
		t.Errorf("Expected task 2 to not be completed, but it succeeded.")
	}
	if msg != "Task 2 is already completed." {
		t.Errorf("Incorrect already completed message '%s'", msg)
	}

	success, msg = tracker.CompleteTask(99)
	if success {
		t.Errorf("Expected task 99 to not be completed, but it succeeded.")
	}
	if msg != "Task with ID 99 not found." {
		t.Errorf("Incorrect not found message. Got: '%s'", msg)
	}

	if tracker.tasks[0].Completed {
		t.Errorf("Task 1 should not be completed")
	}
	if tracker.tasks[2].Completed {
		t.Errorf("Task 3 should not be completed")
	}
}

func TestListTasks(t *testing.T) {
	tracker := NewTaskTracker()

	expected := "Pending Tasks:\nNo pending tasks."
	if tracker.ListTasks() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, tracker.ListTasks())
	}

	tracker.AddTask("Task A")
	expected = "Pending Tasks:\n1: Task A\n"
	if tracker.ListTasks() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, tracker.ListTasks())
	}

	tracker.AddTask("Task B")
	expected = "Pending Tasks:\n1: Task A\n2: Task B\n"
	if tracker.ListTasks() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, tracker.ListTasks())
	}

	tracker.CompleteTask(1)
	expected = "Pending Tasks:\n2: Task B\n"
	if tracker.ListTasks() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, tracker.ListTasks())
	}

	tracker.CompleteTask(2)
	expected = "Pending Tasks:\nNo pending tasks."
	if tracker.ListTasks() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, tracker.ListTasks())
	}
}
