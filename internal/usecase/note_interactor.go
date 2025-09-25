package usecase

import (
	"context"
	"fmt"
	"go_test/internal/domain"
)

// noteInteractor implements NoteUsecase interface
type noteInteractor struct {
	noteRepo NoteRepository
	cache    Cache
}

// NewNoteInteractor creates a new note interactor
func NewNoteInteractor(noteRepo NoteRepository, cache Cache) NoteUsecase {
	return &noteInteractor{
		noteRepo: noteRepo,
		cache:    cache,
	}
}

// CreateNote creates a new note
func (n *noteInteractor) CreateNote(ctx context.Context, title, content string) (*domain.Note, error) {
	note := domain.NewNote(title, content)

	createdNote, err := n.noteRepo.Create(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	// Cache the created note
	cacheKey := fmt.Sprintf("note:%d", createdNote.ID)
	if err := n.cache.Set(ctx, cacheKey, createdNote); err != nil {
		// Log error but don't fail the operation
		// In a real application, you might want to use a logger here
	}

	return createdNote, nil
}

// GetNote retrieves a note by ID
func (n *noteInteractor) GetNote(ctx context.Context, id int64) (*domain.Note, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("note:%d", id)
	if cachedValue, err := n.cache.Get(ctx, cacheKey); err == nil && cachedValue != "" {
		// In a real application, you would deserialize the cached value
		// For simplicity, we'll always fetch from database
	}

	note, err := n.noteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return note, nil
}
