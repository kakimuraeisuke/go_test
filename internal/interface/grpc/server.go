package grpc

import (
	"context"
	"go_test/internal/usecase"
	v1 "go_test/proto/go_test/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server はgRPCサービスを実装します
type server struct {
	v1.UnimplementedGoTestServiceServer
	noteUsecase usecase.NoteUsecase
	pingUsecase usecase.PingUsecase
}

// NewServer は新しいgRPCサーバーを作成します
func NewServer(noteUsecase usecase.NoteUsecase, pingUsecase usecase.PingUsecase) *grpc.Server {
	s := &server{
		noteUsecase: noteUsecase,
		pingUsecase: pingUsecase,
	}

	grpcServer := grpc.NewServer()
	v1.RegisterGoTestServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	
	// ヘルスチェックサービスを登録
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	return grpcServer
}

// Ping はPing RPCメソッドを実装します
func (s *server) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error) {
	mysqlAvailable, redisAvailable, message, err := s.pingUsecase.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ping failed: %v", err)
	}

	return &v1.PingResponse{
		MysqlAvailable: mysqlAvailable,
		RedisAvailable: redisAvailable,
		Message:        message,
	}, nil
}

// CreateNote はCreateNote RPCメソッドを実装します
func (s *server) CreateNote(ctx context.Context, req *v1.CreateNoteRequest) (*v1.CreateNoteResponse, error) {
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	note, err := s.noteUsecase.CreateNote(ctx, req.Title, req.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create note: %v", err)
	}

	return &v1.CreateNoteResponse{
		Id:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: timestamppb.New(note.CreatedAt),
	}, nil
}

// GetNote はGetNote RPCメソッドを実装します
func (s *server) GetNote(ctx context.Context, req *v1.GetNoteRequest) (*v1.GetNoteResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id must be positive")
	}

	note, err := s.noteUsecase.GetNote(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get note: %v", err)
	}

	return &v1.GetNoteResponse{
		Id:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: timestamppb.New(note.CreatedAt),
	}, nil
}
