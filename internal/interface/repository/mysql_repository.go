package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go_test/internal/domain"
	"go_test/internal/usecase"
)

// mysqlRepository implements NoteRepository interface
type mysqlRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new MySQL repository
func NewMySQLRepository(db *sql.DB) usecase.NoteRepository {
	return &mysqlRepository{db: db}
}

// Create creates a new note in the database
func (r *mysqlRepository) Create(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	query := `INSERT INTO notes (title, content) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, note.Title, note.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to insert note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	note.ID = id
	return note, nil
}

// GetByID retrieves a note by ID from the database
func (r *mysqlRepository) GetByID(ctx context.Context, id int64) (*domain.Note, error) {
	query := `SELECT id, title, content, created_at FROM notes WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var note domain.Note
	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to scan note: %w", err)
	}

	return &note, nil
}
