package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go_test/internal/infrastructure/mysql"
	redisInfra "go_test/internal/infrastructure/redis"
	"go_test/internal/interface/cache"
	"go_test/internal/interface/grpc"
	"go_test/internal/interface/repository"
	"go_test/internal/usecase"

	"database/sql"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables
	loadEnv()

	// Initialize MySQL connection
	mysqlConfig := mysql.NewConfig()
	db, err := mysql.Connect(mysqlConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Initialize Redis connection
	redisConfig := redisInfra.NewConfig()
	redisClient := redisInfra.Connect(redisConfig)
	defer redisClient.Close()

	// Test Redis connection
	ctx := context.Background()
	if err := redisInfra.Ping(ctx, redisClient); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	// Initialize repositories and cache
	noteRepo := repository.NewMySQLRepository(db)
	redisCache := cache.NewRedisCache(redisClient)

	// Initialize use cases
	noteUsecase := usecase.NewNoteInteractor(noteRepo, redisCache)
	pingUsecase := usecase.NewPingInteractor(&sqlPinger{db: db}, &redisPinger{client: redisClient})

	// Initialize gRPC server
	grpcServer := grpc.NewServer(noteUsecase, pingUsecase)

	// Start gRPC server
	port := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("Starting gRPC server on port %s", port)

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	grpcServer.GracefulStop()
}

// loadEnv loads environment variables from .env file if it exists
func loadEnv() {
	// In a real application, you might want to use a library like godotenv
	// For simplicity, we'll just check if .env exists and log a message
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Loading environment variables from .env file")
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// sqlPinger implements usecase.SQLPinger interface
type sqlPinger struct {
	db *sql.DB
}

func (p *sqlPinger) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// redisPinger implements usecase.RedisPinger interface
type redisPinger struct {
	client *redis.Client
}

func (p *redisPinger) Ping(ctx context.Context) error {
	_, err := p.client.Ping(ctx).Result()
	return err
}
