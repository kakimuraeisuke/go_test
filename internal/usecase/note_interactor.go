package usecase

import (
	"context"
	"fmt"
	"go_test/internal/domain"
)

// noteInteractor はNoteUsecaseインターフェースを実装します
type noteInteractor struct {
	noteRepo NoteRepository
	cache    Cache
}

// NewNoteInteractor は新しいノートインタラクターを作成します
func NewNoteInteractor(noteRepo NoteRepository, cache Cache) NoteUsecase {
	return &noteInteractor{
		noteRepo: noteRepo,
		cache:    cache,
	}
}

// CreateNote は新しいノートを作成します
func (n *noteInteractor) CreateNote(ctx context.Context, title, content string) (*domain.Note, error) {
	note := domain.NewNote(title, content)
	
	createdNote, err := n.noteRepo.Create(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	// 作成されたノートをキャッシュに保存
	cacheKey := fmt.Sprintf("note:%d", createdNote.ID)
	if err := n.cache.Set(ctx, cacheKey, createdNote); err != nil {
		// エラーをログに記録しますが、操作は失敗させません
		// 実際のアプリケーションでは、ロガーを使用することを推奨します
	}

	return createdNote, nil
}

// GetNote はIDでノートを取得します
func (n *noteInteractor) GetNote(ctx context.Context, id int64) (*domain.Note, error) {
	// まずキャッシュから取得を試行
	cacheKey := fmt.Sprintf("note:%d", id)
	if cachedValue, err := n.cache.Get(ctx, cacheKey); err == nil && cachedValue != "" {
		// 実際のアプリケーションでは、キャッシュされた値をデシリアライズします
		// 簡略化のため、常にデータベースから取得します
	}

	note, err := n.noteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return note, nil
}
