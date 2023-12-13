package auth

import (
	pb "auth_service/internal/api"
	"auth_service/internal/configuration"
	"context"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

// Auth defines the interface to interact with jwt tokens. Includes generating, validation, refreshing and data extraction.
type Auth interface {
	GenerateToken(ctx context.Context, req *pb.GenerateTokenRequest) (*pb.TokenResponse, error)
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.TokenValidationResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.TokenResponse, error)
	ExtractTokenData(ctx context.Context, req *pb.ExtractTokenDataRequest) (*pb.TokenDataResponse, error)
}

// JwtAuthService is the implementation of Auth .
type JwtAuthService struct {
	pb.UnimplementedAuthServiceServer
	auth       Auth
	publicKey  rsa.PublicKey
	privateKey rsa.PrivateKey

	signingMethod jwt.SigningMethod

	log *zap.SugaredLogger
}

func NewJWTAuthService(config *configuration.AuthServiceConfig, logger *zap.SugaredLogger) (*JwtAuthService, error) {
	// Loading the raw private key from path defined in config
	privateKeyData, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	// Decoding of private key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	// Loading the raw public key from path defined in config
	publicKeyData, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	// Decoding of public key
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &JwtAuthService{
		log:           logger,
		privateKey:    *privateKey,
		publicKey:     *publicKey,
		signingMethod: jwt.SigningMethodRS256,
	}, nil
}

func (s *JwtAuthService) GenerateToken(ctx context.Context, in *pb.GenerateTokenRequest) (*pb.TokenResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "not implemented yet")
}

func (s *JwtAuthService) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.TokenValidationResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "not implemented yet")
}

func (s *JwtAuthService) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "not implemented yet")
}

func (s *JwtAuthService) ExtractTokenData(ctx context.Context, in *pb.ExtractTokenDataRequest) (*pb.TokenDataResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "not implemented yet")
}
