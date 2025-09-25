# Go Test Service

クリーンアーキテクチャに準拠したGo + gRPC + MySQL + Redisのマイクロサービスです。

## アーキテクチャ

このプロジェクトはクリーンアーキテクチャの原則に従って設計されています：

- **Domain層**: ビジネスロジックの中核となるエンティティ
- **UseCase層**: アプリケーション固有のビジネスルール
- **Interface層**: 外部との接続点（gRPC、リポジトリ、キャッシュ）
- **Infrastructure層**: 外部システムとの接続（MySQL、Redis）

## サービス仕様

### gRPCサービス
- サービス名: `go_test.v1.GoTestService`
- メソッド:
  - `Ping`: MySQL/Redisの到達性を確認
  - `CreateNote`: ノートを作成
  - `GetNote`: ノートを取得

### データモデル
- データベース: `go_test`
- テーブル: `notes`
  - `id`: BIGINT AUTO_INCREMENT PRIMARY KEY
  - `title`: VARCHAR(255)
  - `content`: TEXT
  - `created_at`: TIMESTAMP DEFAULT CURRENT_TIMESTAMP

## セットアップ

### 1. 環境変数の設定

```bash
cp env.example .env
```

必要に応じて`.env`ファイルの値を調整してください。

### 2. アプリケーションの起動

```bash
docker compose up --build
```

### 3. 動作確認

1. **grpcuiでの確認**
   - ブラウザで `http://localhost:8080` にアクセス
   - gRPCサービスの一覧が表示されます

2. **Pingテスト**
   - `Ping`メソッドを実行
   - MySQLとRedisの接続状況が確認できます

3. **Note操作テスト**
   - `CreateNote`でノートを作成
   - `GetNote`で作成したノートを取得

## 開発

### protobufファイルの生成

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/go_test/v1/go_test.proto
```

### ローカル開発

1. MySQLとRedisを起動
```bash
docker compose up mysql redis
```

2. アプリケーションを起動
```bash
go run cmd/server/main.go
```

## トラブルシューティング

### よくある問題

1. **protocコマンドが見つからない**
   - Protocol Buffersをインストールしてください
   - macOS: `brew install protobuf`
   - Ubuntu: `apt-get install protobuf-compiler`

2. **MySQL接続エラー**
   - MySQLコンテナが正常に起動しているか確認
   - 環境変数が正しく設定されているか確認

3. **Redis接続エラー**
   - Redisコンテナが正常に起動しているか確認
   - ネットワーク設定を確認

4. **grpcuiにアクセスできない**
   - ポート8080が使用可能か確認
   - ファイアウォール設定を確認

### ヘルスチェック

```bash
# MySQLのヘルスチェック
docker compose exec mysql mysqladmin ping -h localhost -u root -prootpassword

# Redisのヘルスチェック
docker compose exec redis redis-cli ping

# gRPCサービスのヘルスチェック
docker compose exec goapp grpc-health-probe -addr=:50051
```

## ディレクトリ構成

```
.
├─ cmd/server/main.go              # アプリケーションエントリーポイント
├─ internal/
│  ├─ domain/note.go              # ドメインエンティティ
│  ├─ usecase/                    # ユースケース層
│  │  ├─ ports.go                # インターフェース定義
│  │  ├─ note_interactor.go      # ノートユースケース実装
│  │  └─ ping_interactor.go      # ピングユースケース実装
│  ├─ interface/                 # インターフェース層
│  │  ├─ grpc/server.go          # gRPCサーバー
│  │  ├─ repository/mysql_repository.go  # MySQLリポジトリ
│  │  └─ cache/redis_cache.go    # Redisキャッシュ
│  └─ infrastructure/            # インフラストラクチャ層
│     ├─ mysql/conn.go           # MySQL接続
│     └─ redis/conn.go           # Redis接続
├─ proto/go_test/v1/go_test.proto # protobuf定義
├─ mysql/                        # MySQL設定
│  ├─ conf.d/my.cnf
│  └─ initdb.d/001_create_db_and_table.sql
├─ Dockerfile
├─ docker-compose.yml
├─ go.mod
├─ env.example
└─ README.md
```

## 依存関係

- Go 1.22+
- Protocol Buffers
- MySQL 8.4
- Redis
- Docker & Docker Compose
