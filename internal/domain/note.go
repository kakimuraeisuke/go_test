package domain

import "time"

// Note represents a note entity in the domain layer
type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// NewNote creates a new Note instance
func NewNote(title, content string) *Note {
	return &Note{
		Title:   title,
		Content: content,
	}
}
