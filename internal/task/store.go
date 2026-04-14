package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// JSONFileStore provides thread-safe persistence for tasks to a JSON file.
// Uses sync.Mutex for concurrent access protection (DATA-01).
type JSONFileStore struct {
	mu   sync.Mutex
	path string
}

// NewJSONFileStore creates a new JSONFileStore.
// If path is empty, defaults to ~/.todo.json via os.UserHomeDir().
func NewJSONFileStore(path string) (*JSONFileStore, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".todo.json")
	}
	return &JSONFileStore{path: path}, nil
}

// Load reads and returns all tasks from the JSON file.
// Returns an empty Task slice if the file doesn't exist (DATA-02).
func (s *JSONFileStore) Load() ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

// loadLocked reads tasks assuming mutex is already held.
func (s *JSONFileStore) loadLocked() ([]Task, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // DATA-02: first run returns empty
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Task{}, nil
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("corrupted todo file: %w", err)
	}
	return tasks, nil
}

// Save writes the tasks slice to the JSON file atomically (DATA-03).
// Uses CreateTemp + Rename pattern for atomic writes.
func (s *JSONFileStore) Save(tasks []Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked(tasks)
}

// saveLocked writes tasks assuming mutex is already held.
func (s *JSONFileStore) saveLocked(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	// CreateTemp creates a FILE (not a directory) in the parent directory.
	// Use filepath.Dir(s.path) to put temp file in same directory as target.
	f, err := os.CreateTemp(filepath.Dir(s.path), "todo-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	f.Close() // Close file handle before WriteFile

	defer os.Remove(tmp) // safe if file already moved

	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Add appends a new task to the store and persists immediately.
func (s *JSONFileStore) Add(task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks, err := s.loadLocked()
	if err != nil {
		return err
	}
	tasks = append(tasks, *task)
	return s.saveLocked(tasks)
}

// List returns all tasks, optionally filtered by category (CAT-02).
// If category is empty string, returns all tasks.
func (s *JSONFileStore) List(category string) ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks, err := s.loadLocked()
	if err != nil {
		return nil, err
	}
	if category == "" {
		return tasks, nil
	}
	filtered := make([]Task, 0)
	for _, t := range tasks {
		if t.Category == category {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

// MarkDone sets the Completed field to true for the task with the given ID.
// Returns an error if the task is not found.
func (s *JSONFileStore) MarkDone(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks, err := s.loadLocked()
	if err != nil {
		return err
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Completed = true
			return s.saveLocked(tasks)
		}
	}
	return fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)
}

// Delete removes the task with the given ID from the store.
// Returns an error if the task is not found.
func (s *JSONFileStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks, err := s.loadLocked()
	if err != nil {
		return err
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return s.saveLocked(tasks)
		}
	}
	return fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)
}

// ErrTaskNotFound is returned when a task is not found (placeholder for Phase 2).
var ErrTaskNotFound = errors.New("task not found")
