package usecase

import (
	"context"
	"fmt"
)

// pingInteractor はPingUsecaseインターフェースを実装します
type pingInteractor struct {
	sqlPinger   SQLPinger
	redisPinger RedisPinger
}

// NewPingInteractor は新しいピングインタラクターを作成します
func NewPingInteractor(sqlPinger SQLPinger, redisPinger RedisPinger) PingUsecase {
	return &pingInteractor{
		sqlPinger:   sqlPinger,
		redisPinger: redisPinger,
	}
}

// Ping はMySQLとRedisの可用性をチェックします
func (p *pingInteractor) Ping(ctx context.Context) (mysqlAvailable, redisAvailable bool, message string, err error) {
	// MySQLの可用性をチェック
	mysqlErr := p.sqlPinger.Ping(ctx)
	mysqlAvailable = mysqlErr == nil

	// Redisの可用性をチェック
	redisErr := p.redisPinger.Ping(ctx)
	redisAvailable = redisErr == nil

	// メッセージを構築
	if mysqlAvailable && redisAvailable {
		message = "All services are available"
	} else if mysqlAvailable {
		message = fmt.Sprintf("MySQL is available, Redis is not: %v", redisErr)
	} else if redisAvailable {
		message = fmt.Sprintf("Redis is available, MySQL is not: %v", mysqlErr)
	} else {
		message = fmt.Sprintf("Both services are unavailable. MySQL: %v, Redis: %v", mysqlErr, redisErr)
	}

	return mysqlAvailable, redisAvailable, message, nil
}
