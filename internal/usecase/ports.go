package usecase

import (
	"context"
	"go_test/internal/domain"
)

// NoteUsecase はノートユースケースのインターフェースを定義します
type NoteUsecase interface {
	CreateNote(ctx context.Context, title, content string) (*domain.Note, error)
	GetNote(ctx context.Context, id int64) (*domain.Note, error)
}

// PingUsecase はピングユースケースのインターフェースを定義します
type PingUsecase interface {
	Ping(ctx context.Context) (mysqlAvailable, redisAvailable bool, message string, err error)
}

// NoteRepository はノートリポジトリのインターフェースを定義します
type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) (*domain.Note, error)
	GetByID(ctx context.Context, id int64) (*domain.Note, error)
}

// Cache はキャッシュ操作のインターフェースを定義します
type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

// SQLPinger はSQLピング操作のインターフェースを定義します
type SQLPinger interface {
	Ping(ctx context.Context) error
}

// RedisPinger はRedisピング操作のインターフェースを定義します
type RedisPinger interface {
	Ping(ctx context.Context) error
}
