package usecase

import (
	"context"
	"fmt"
)

// pingInteractor implements PingUsecase interface
type pingInteractor struct {
	sqlPinger   SQLPinger
	redisPinger RedisPinger
}

// NewPingInteractor creates a new ping interactor
func NewPingInteractor(sqlPinger SQLPinger, redisPinger RedisPinger) PingUsecase {
	return &pingInteractor{
		sqlPinger:   sqlPinger,
		redisPinger: redisPinger,
	}
}

// Ping checks the availability of MySQL and Redis
func (p *pingInteractor) Ping(ctx context.Context) (mysqlAvailable, redisAvailable bool, message string, err error) {
	// Check MySQL availability
	mysqlErr := p.sqlPinger.Ping(ctx)
	mysqlAvailable = mysqlErr == nil

	// Check Redis availability
	redisErr := p.redisPinger.Ping(ctx)
	redisAvailable = redisErr == nil

	// Build message
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
