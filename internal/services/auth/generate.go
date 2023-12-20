package auth

import (
	pb "auth_service/internal/api"
	"auth_service/internal/domain/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func generateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *JwtAuthService) GenerateToken(ctx context.Context, in *pb.GenerateTokenRequest) (*pb.TokenResponse, error) {
	s.log.Debug("Generating token")

	// Check if user_id is empty
	if in.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is empty")
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := s.generateTokens(ctx, in.UserId, in.UserData)
	if err != nil {
		return nil, err
	}

	s.log.Debug("Tokens generated successfully")

	// Sending the response
	return &pb.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateTokens creates new access and refresh tokens
func (s *JwtAuthService) generateTokens(ctx context.Context, userId string, userData *pb.UserData) (string, string, error) {
	// Adding user_id to claims
	claims := jwt.MapClaims{
		"user_id": userId,
	}

	// Adding additional user data
	for key, value := range userData.AdditionalData {
		claims[key] = value
	}

	// Adding exp field to the claims
	jwtTTL := s.tokenTTL
	jwtExpiresAt := time.Now().Add(jwtTTL)
	claims["exp"] = jwtExpiresAt.Unix()

	// Creating the access token
	token := jwt.NewWithClaims(s.signingMethod, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		s.log.Errorf("Failed to sign access token: %v", err)
		return "", "", status.Errorf(codes.Internal, "failed to sign access token")
	}

	// Generating the refresh token
	refreshToken, err := generateRandomString(36)
	if err != nil {
		s.log.Errorf("Failed to generate refresh token: %v", err)
		return "", "", status.Errorf(codes.Internal, "failed to generate refresh token")
	}

	// Create a Redis entry for the refresh session
	expiresAt := time.Now().Add(s.refreshTokenTTL)
	redisEntry := models.RefreshSession{
		UserID:       userId,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}

	// Inserting the entry to the redis
	err = s.redisRepository.InsertRefreshSession(ctx, &redisEntry)
	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to insert refresh token to the redis storage")
	}

	return signedToken, refreshToken, nil
}
