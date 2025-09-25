package domain

import "time"

// Note はドメイン層のノートエンティティを表します
type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// NewNote は新しいNoteインスタンスを作成します
func NewNote(title, content string) *Note {
	return &Note{
		Title:   title,
		Content: content,
	}
}
