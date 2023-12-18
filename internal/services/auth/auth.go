package auth

import (
	pb "auth_service/internal/api"
	"auth_service/internal/configuration"
	"auth_service/internal/repository"
	"context"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

// Auth defines the interface to interact with jwt tokens. Includes generating, validation, refreshing and data extraction.
// Also defines Shutdown method to gracefully shutdown redis connection.
type Auth interface {
	GenerateToken(ctx context.Context, req *pb.GenerateTokenRequest) (*pb.TokenResponse, error)
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.TokenValidationResponse, error)
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.TokenResponse, error)
	ExtractTokenData(ctx context.Context, req *pb.ExtractTokenDataRequest) (*pb.TokenDataResponse, error)
	Shutdown() error
}

// JwtAuthService is the implementation of Auth .
type JwtAuthService struct {
	pb.UnimplementedAuthServiceServer
	auth       Auth
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	signingMethod jwt.SigningMethod

	refreshTokenTTL time.Duration
	tokenTTL        time.Duration

	redisRepository repository.RefreshSessionsRepository
	log             *zap.SugaredLogger
}

func NewJWTAuthService(config *configuration.AuthServiceConfig, logger *zap.SugaredLogger, redisRepository repository.RefreshSessionsRepository, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *JwtAuthService {

	// Converting refresh token TTL value from config to time.Duration
	refreshTokenTTL := time.Hour * 24 * time.Duration(config.RefreshTokenLifeTime)

	// Converting standard(not a refresh one) token TTL value from config to time.Duration
	tokenTTL := time.Minute * time.Duration(config.RefreshTokenLifeTime)

	return &JwtAuthService{
		privateKey:      privateKey,
		publicKey:       publicKey,
		signingMethod:   jwt.SigningMethodRS256,
		redisRepository: redisRepository,
		refreshTokenTTL: refreshTokenTTL,
		tokenTTL:        tokenTTL,
		log:             logger,
	}
}

func (s *JwtAuthService) Shutdown() error {

	// Disconnecting from redis
	if err := s.redisRepository.Disconnect(); err != nil {
		s.log.Errorf("Failed to close Redis client: %v", err)
		return err
	}

	s.log.Info("Successfully shut down the service")
	return nil
}
