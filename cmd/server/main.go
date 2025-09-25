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
	// 環境変数を読み込み
	loadEnv()

	// MySQL接続を初期化
	mysqlConfig := mysql.NewConfig()
	db, err := mysql.Connect(mysqlConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Redis接続を初期化
	redisConfig := redisInfra.NewConfig()
	redisClient := redisInfra.Connect(redisConfig)
	defer redisClient.Close()

	// Redis接続をテスト
	ctx := context.Background()
	if err := redisInfra.Ping(ctx, redisClient); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	// リポジトリとキャッシュを初期化
	noteRepo := repository.NewMySQLRepository(db)
	redisCache := cache.NewRedisCache(redisClient)

	// ユースケースを初期化
	noteUsecase := usecase.NewNoteInteractor(noteRepo, redisCache)
	pingUsecase := usecase.NewPingInteractor(&sqlPinger{db: db}, &redisPinger{client: redisClient})

	// gRPCサーバーを初期化
	grpcServer := grpc.NewServer(noteUsecase, pingUsecase)

	// gRPCサーバーを開始
	port := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("Starting gRPC server on port %s", port)

	// サーバーをgoroutineで開始
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// 割り込みシグナルを待ってサーバーを正常にシャットダウン
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 正常なシャットダウン
	grpcServer.GracefulStop()
}

// loadEnv は.envファイルが存在する場合に環境変数を読み込みます
func loadEnv() {
	// 実際のアプリケーションでは、godotenvのようなライブラリを使用することを推奨します
	// 簡略化のため、.envファイルの存在をチェックしてメッセージをログに記録します
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Loading environment variables from .env file")
	}
}

// getEnv はデフォルト値付きで環境変数を取得します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// sqlPinger はusecase.SQLPingerインターフェースを実装します
type sqlPinger struct {
	db *sql.DB
}

func (p *sqlPinger) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// redisPinger はusecase.RedisPingerインターフェースを実装します
type redisPinger struct {
	client *redis.Client
}

func (p *redisPinger) Ping(ctx context.Context) error {
	_, err := p.client.Ping(ctx).Result()
	return err
}
