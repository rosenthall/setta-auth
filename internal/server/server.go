package server

import (
	pb "auth_service/internal/api"
	"auth_service/internal/configuration"
	auth "auth_service/internal/services"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

// AuthServer contains all the state our service need.
type AuthServer struct {
	grpcServer  *grpc.Server
	config      *configuration.AuthServiceConfig
	logger      *zap.SugaredLogger
	appProvider auth.Auth
	pb.AuthServiceServer
}

func (s *AuthServer) GenerateToken(ctx context.Context, request *pb.GenerateTokenRequest) (*pb.TokenResponse, error) {
	return s.appProvider.GenerateToken(ctx, request)
}

func (s *AuthServer) ValidateToken(ctx context.Context, request *pb.ValidateTokenRequest) (*pb.TokenValidationResponse, error) {
	return s.appProvider.ValidateToken(ctx, request)
}

func (s *AuthServer) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	return s.appProvider.RefreshToken(ctx, request)
}

func (s *AuthServer) ExtractTokenData(ctx context.Context, request *pb.ExtractTokenDataRequest) (*pb.TokenDataResponse, error) {
	return s.appProvider.ExtractTokenData(ctx, request)
}

func (s *AuthServer) mustEmbedUnimplementedAuthServiceServer() {
	//TODO implement me
	panic("implement me")
}

// NewServer creates an instance of AuthServer.
func NewServer(config *configuration.AuthServiceConfig, logger *zap.SugaredLogger, appProvider auth.Auth) *AuthServer {
	grpcServer := grpc.NewServer()
	return &AuthServer{
		grpcServer:  grpcServer,
		config:      config,
		logger:      logger,
		appProvider: appProvider,
	}
}

func (s *AuthServer) Run() error {
	address := fmt.Sprintf(":%d", s.config.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Registration of our service
	pb.RegisterAuthServiceServer(s.grpcServer, s)

	// Starting server
	s.logger.Info("Starting gRPC server on 127.0.0.1", address)
	return s.grpcServer.Serve(lis)
}
