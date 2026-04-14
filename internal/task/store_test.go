package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestNewJSONFileStore(t *testing.T) {
	// Test with custom path
	store, err := NewJSONFileStore("/custom/path/todo.json")
	if err != nil {
		t.Fatalf("NewJSONFileStore failed: %v", err)
	}
	if store == nil {
		t.Fatal("NewJSONFileStore returned nil store")
	}
	if store.path != "/custom/path/todo.json" {
		t.Errorf("expected path /custom/path/todo.json, got %s", store.path)
	}

	// Test with empty path defaults to ~/.todo.json
	store, err = NewJSONFileStore("")
	if err != nil {
		t.Fatalf("NewJSONFileStore with empty path failed: %v", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir failed: %v", err)
	}
	expected := filepath.Join(home, ".todo.json")
	if store.path != expected {
		t.Errorf("expected default path %s, got %s", expected, store.path)
	}
}

func TestLoadNotExist(t *testing.T) {
	// Create store with non-existent path
	store, err := NewJSONFileStore("/nonexistent/path/todo.json")
	if err != nil {
		t.Fatalf("NewJSONFileStore failed: %v", err)
	}

	// Load should return empty slice, not error (DATA-02)
	tasks, err := store.Load()
	if err != nil {
		t.Fatalf("Load on non-existent file returned error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty slice, got %d tasks", len(tasks))
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temp directory for test file
	tmpDir, err := os.MkdirTemp("", "todo-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "test-todo.json")
	store, err := NewJSONFileStore(testPath)
	if err != nil {
		t.Fatalf("NewJSONFileStore failed: %v", err)
	}

	// Create and save some tasks
	task1, _ := NewTask("Task 1", "work")
	task2, _ := NewTask("Task 2", "home")

	err = store.Save([]Task{*task1, *task2})
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load should return the saved tasks
	tasks, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "Task 1" {
		t.Errorf("expected task 1 title 'Task 1', got '%s'", tasks[0].Title)
	}
	if tasks[1].Title != "Task 2" {
		t.Errorf("expected task 2 title 'Task 2', got '%s'", tasks[1].Title)
	}
}

func TestAdd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "todo-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "test-todo.json")
	store, err := NewJSONFileStore(testPath)
	if err != nil {
		t.Fatalf("NewJSONFileStore failed: %v", err)
	}

	// Add a task
	task, err := NewTask("New Task", "testing")
	if err != nil {
		t.Fatalf("NewTask failed: %v", err)
	}

	err = store.Add(task)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Verify task persists after reload
	tasks, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task after Add, got %d", len(tasks))
	}
	if tasks[0].Title != "New Task" {
		t.Errorf("expected title 'New Task', got '%s'", tasks[0].Title)
	}

	// Add another task
	task2, _ := NewTask("Second Task", "testing")
	err = store.Add(task2)
	if err != nil {
		t.Fatalf("Add second task failed: %v", err)
	}

	tasks, err = store.Load()
	if err != nil {
		t.Fatalf("Load after second Add failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks after second Add, got %d", len(tasks))
	}
}

func TestList(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "todo-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "test-todo.json")
	store, err := NewJSONFileStore(testPath)
	if err != nil {
		t.Fatalf("NewJSONFileStore failed: %v", err)
	}

	// Add tasks with different categories
	task1, _ := NewTask("Task 1", "work")
	task2, _ := NewTask("Task 2", "home")
	task3, _ := NewTask("Task 3", "work")

	for _, task := range []*Task{task1, task2, task3} {
		if err := store.Add(task); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	// List all tasks (empty category filter)
	allTasks, err := store.List("")
	if err != nil {
		t.Fatalf("List('') failed: %v", err)
	}
	if len(allTasks) != 3 {
		t.Errorf("expected 3 tasks with empty filter, got %d", len(allTasks))
	}

	// List filtered by "work"
	workTasks, err := store.List("work")
	if err != nil {
		t.Fatalf("List('work') failed: %v", err)
	}
	if len(workTasks) != 2 {
		t.Errorf("expected 2 work tasks, got %d", len(workTasks))
	}
	for _, task := range workTasks {
		if task.Category != "work" {
			t.Errorf("expected category 'work', got '%s'", task.Category)
		}
	}

	// List filtered by "home"
	homeTasks, err := store.List("home")
	if err != nil {
		t.Fatalf("List('home') failed: %v", err)
	}
	if len(homeTasks) != 1 {
		t.Errorf("expected 1 home task, got %d", len(homeTasks))
	}

	// List filtered by non-existent category
	otherTasks, err := store.List("nonexistent")
	if err != nil {
		t.Fatalf("List('nonexistent') failed: %v", err)
	}
	if len(otherTasks) != 0 {
		t.Errorf("expected 0 tasks for nonexistent category, got %d", len(otherTasks))
	}
}

func TestMarkDone(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "todo-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "test-todo.json")
	store, err := NewJSONFileStore(testPath)
	if err != nil {
		t.Fatalf("NewJSONFileStore failed: %v", err)
	}

	// Add tasks
	task1, _ := NewTask("Task 1", "work")
	task2, _ := NewTask("Task 2", "home")

	for _, task := range []*Task{task1, task2} {
		if err := store.Add(task); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	// Mark the first task done
	err = store.MarkDone(task1.ID)
	if err != nil {
		t.Fatalf("MarkDone failed: %v", err)
	}

	// Verify task is marked done
	tasks, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	for _, task := range tasks {
		if task.ID == task1.ID {
			if !task.Completed {
				t.Error("expected task to be marked completed")
			}
		} else {
			if task.Completed {
				t.Error("task that wasn't marked done is marked completed")
			}
		}
	}

	// MarkDone on non-existent ID should error
	err = store.MarkDone("nonexistent-id")
	if err == nil {
		t.Error("expected error for MarkDone on non-existent ID")
	}
}

func TestDelete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "todo-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "test-todo.json")
	store, err := NewJSONFileStore(testPath)
	if err != nil {
		t.Fatalf("NewJSONFileStore failed: %v", err)
	}

	// Add tasks
	task1, _ := NewTask("Task 1", "work")
	task2, _ := NewTask("Task 2", "home")

	for _, task := range []*Task{task1, task2} {
		if err := store.Add(task); err != nil {
			t.Fatalf("Add failed: %v", err)
		}
	}

	// Delete the first task
	err = store.Delete(task1.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify only one task remains
	tasks, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task after Delete, got %d", len(tasks))
	}
	if tasks[0].ID != task2.ID {
		t.Errorf("expected remaining task to be task2, got ID %s", tasks[0].ID)
	}

	// Delete on non-existent ID should error
	err = store.Delete("nonexistent-id")
	if err == nil {
		t.Error("expected error for Delete on non-existent ID")
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Create temp store
	tmpDir, err := os.MkdirTemp("", "todo-concurrent-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "concurrent.json")
	store, err := NewJSONFileStore(storePath)
	if err != nil {
		t.Fatal(err)
	}

	// Launch concurrent operations
	var wg sync.WaitGroup
	numGoroutines := 10
	numOpsPerGoroutine := 20

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOpsPerGoroutine; j++ {
				task, err := NewTask(fmt.Sprintf("Task-%d-%d", goroutineID, j), "")
				if err != nil {
					t.Errorf("NewTask error: %v", err)
					return
				}
				if err := store.Add(task); err != nil {
					t.Errorf("Add error: %v", err)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify all tasks were added (no data loss)
	allTasks, err := store.List("")
	if err != nil {
		t.Fatalf("List error: %v", err)
	}

	expectedCount := numGoroutines * numOpsPerGoroutine
	if len(allTasks) != expectedCount {
		t.Errorf("Expected %d tasks, got %d (data loss detected)", expectedCount, len(allTasks))
	}

	// Verify file is valid JSON (no corruption)
	data, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	var verifyTasks []Task
	if err := json.Unmarshal(data, &verifyTasks); err != nil {
		t.Fatalf("JSON corruption detected: %v", err)
	}

	if len(verifyTasks) != expectedCount {
		t.Errorf("File JSON has %d tasks but store has %d (file/store mismatch)", len(verifyTasks), expectedCount)
	}
}

func TestAtomicWrite(t *testing.T) {
	// Create temp store
	tmpDir, err := os.MkdirTemp("", "todo-atomic-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "atomic.json")
	store, err := NewJSONFileStore(storePath)
	if err != nil {
		t.Fatal(err)
	}

	// Add initial tasks
	for i := 0; i < 5; i++ {
		task, _ := NewTask(fmt.Sprintf("Initial-%d", i), "")
		store.Add(task)
	}

	// Verify initial state is valid JSON
	data1, _ := os.ReadFile(storePath)
	var tasks1 []Task
	if err := json.Unmarshal(data1, &tasks1); err != nil {
		t.Fatalf("Initial file corrupted: %v", err)
	}

	// Add more tasks
	for i := 0; i < 5; i++ {
		task, _ := NewTask(fmt.Sprintf("After-%d", i), "")
		store.Add(task)
	}

	// Verify new state is valid JSON
	data2, _ := os.ReadFile(storePath)
	var tasks2 []Task
	if err := json.Unmarshal(data2, &tasks2); err != nil {
		t.Fatalf("After write file corrupted: %v", err)
	}

	if len(tasks2) != 10 {
		t.Errorf("Expected 10 tasks after atomic write, got %d", len(tasks2))
	}

	// Verify file has proper JSON structure (array starts with [)
	if len(data2) == 0 || data2[0] != '[' {
		t.Errorf("File does not contain valid JSON array")
	}

	// Verify file ends with ] (complete JSON)
	trimmed := bytes.TrimSpace(data2)
	if len(trimmed) == 0 || trimmed[len(trimmed)-1] != ']' {
		t.Errorf("JSON array not properly closed")
	}
}
