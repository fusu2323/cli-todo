package task

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// generateUUID creates a 32-character hex string using cryptographically
// secure random bytes. Returns an error if the random source fails.
func generateUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil // 32-char hex string
}

// Task represents a single todo item.
type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category,omitempty"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTask creates a new task with a generated UUID and current timestamp.
// The task starts with Completed set to false.
func NewTask(title, category string) (*Task, error) {
	id, err := generateUUID()
	if err != nil {
		return nil, err
	}
	return &Task{
		ID:        id,
		Title:     title,
		Category:  category,
		Completed: false,
		CreatedAt: time.Now(),
	}, nil
}
