package server

import (
	pb "auth_service/internal/api"
	"auth_service/internal/configuration"
	"auth_service/internal/server/middleware"
	"auth_service/internal/services/auth"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
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

// NewServer creates an instance of AuthServer.
func NewServer(config *configuration.AuthServiceConfig, logger *zap.SugaredLogger, appProvider auth.Auth) *AuthServer {

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryLoggingInterceptor(logger)),
	)

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

	// Starting server in a separate goroutine
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Errorf("Failed to serve: %v", err)
		}
	}()
	s.logger.Info("Starting gRPC server on 127.0.0.1", address)

	// Setting up channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Blocking until a signal is received
	<-quit
	s.logger.Info("Shutting down gRPC server...")

	// Gracefully stopping the server
	s.grpcServer.GracefulStop()

	// Calling the Shutdown method of the auth service
	if err := s.appProvider.Shutdown(); err != nil {
		s.logger.Errorf("Failed to shutdown auth service: %v", err)
		return err
	}

	s.logger.Info("Server gracefully stopped")
	return nil
}

func (s *AuthServer) StopGracefully() {
	s.logger.Info("Gracefully stopping gRPC server")
	err := s.appProvider.Shutdown()
	if err != nil {
		return
	}

	s.grpcServer.GracefulStop()
}
