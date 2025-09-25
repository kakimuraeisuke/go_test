package usecase

import (
	"context"
	"go_test/internal/domain"
)

// NoteUsecase defines the interface for note use cases
type NoteUsecase interface {
	CreateNote(ctx context.Context, title, content string) (*domain.Note, error)
	GetNote(ctx context.Context, id int64) (*domain.Note, error)
}

// PingUsecase defines the interface for ping use cases
type PingUsecase interface {
	Ping(ctx context.Context) (mysqlAvailable, redisAvailable bool, message string, err error)
}

// NoteRepository defines the interface for note repository
type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) (*domain.Note, error)
	GetByID(ctx context.Context, id int64) (*domain.Note, error)
}

// Cache defines the interface for cache operations
type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

// SQLPinger defines the interface for SQL ping operations
type SQLPinger interface {
	Ping(ctx context.Context) error
}

// RedisPinger defines the interface for Redis ping operations
type RedisPinger interface {
	Ping(ctx context.Context) error
}
