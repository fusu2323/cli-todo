package task

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	title := "Buy groceries"
	category := "personal"

	task, err := NewTask(title, category)
	if err != nil {
		t.Fatalf("NewTask returned error: %v", err)
	}

	// ID should be 32-character hex string
	if len(task.ID) != 32 {
		t.Errorf("Expected ID length 32, got %d", len(task.ID))
	}
	hexPattern := regexp.MustCompile(`^[a-f0-9]{32}$`)
	if !hexPattern.MatchString(task.ID) {
		t.Errorf("ID %q is not a valid 32-char hex string", task.ID)
	}

	// Title should be set correctly
	if task.Title != title {
		t.Errorf("Expected Title %q, got %q", title, task.Title)
	}

	// Category should be set correctly
	if task.Category != category {
		t.Errorf("Expected Category %q, got %q", category, task.Category)
	}

	// Completed should be false
	if task.Completed != false {
		t.Errorf("Expected Completed false, got %v", task.Completed)
	}

	// CreatedAt should be within last 5 seconds
	now := time.Now()
	if task.CreatedAt.Before(now.Add(-5*time.Second)) || task.CreatedAt.After(now) {
		t.Errorf("CreatedAt %v is not within expected range", task.CreatedAt)
	}
}

func TestGenerateUUID(t *testing.T) {
	uuid1, err := generateUUID()
	if err != nil {
		t.Fatalf("generateUUID returned error: %v", err)
	}

	uuid2, err := generateUUID()
	if err != nil {
		t.Fatalf("generateUUID returned error on second call: %v", err)
	}

	// Two calls should produce different UUIDs
	if uuid1 == uuid2 {
		t.Error("Two calls to generateUUID produced the same UUID")
	}

	// Output should be 32-character hex
	hexPattern := regexp.MustCompile(`^[a-f0-9]{32}$`)
	if !hexPattern.MatchString(uuid1) {
		t.Errorf("UUID1 %q is not a valid 32-char hex string", uuid1)
	}
	if !hexPattern.MatchString(uuid2) {
		t.Errorf("UUID2 %q is not a valid 32-char hex string", uuid2)
	}
}

func TestTaskJSON(t *testing.T) {
	task := &Task{
		ID:        "a3f1b2c8d4e5f6789012345678901234",
		Title:     "Test task",
		Category:  "work",
		Completed: true,
		CreatedAt: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	// Unmarshal back and verify field names
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	expectedFields := []string{"id", "title", "category", "completed", "created_at"}
	for _, field := range expectedFields {
		if _, ok := decoded[field]; !ok {
			t.Errorf("JSON missing expected field %q", field)
		}
	}

	// Verify values
	if decoded["id"] != "a3f1b2c8d4e5f6789012345678901234" {
		t.Errorf("Expected id 'a3f1b2c8d4e5f6789012345678901234', got %v", decoded["id"])
	}
	if decoded["title"] != "Test task" {
		t.Errorf("Expected title 'Test task', got %v", decoded["title"])
	}
	if decoded["category"] != "work" {
		t.Errorf("Expected category 'work', got %v", decoded["category"])
	}
	if decoded["completed"] != true {
		t.Errorf("Expected completed true, got %v", decoded["completed"])
	}
	if decoded["created_at"] != "2026-01-15T10:30:00Z" {
		t.Errorf("Expected created_at '2026-01-15T10:30:00Z', got %v", decoded["created_at"])
	}
}

func TestTaskJSONOmitempty(t *testing.T) {
	// Task with empty category should omit category in JSON
	task := &Task{
		ID:        "a3f1b2c8d4e5f6789012345678901234",
		Title:     "No category",
		Category:  "",
		Completed: false,
		CreatedAt: time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	// Empty category should be omitted
	if _, ok := decoded["category"]; ok {
		t.Error("Empty category should be omitted from JSON due to omitempty")
	}
}
