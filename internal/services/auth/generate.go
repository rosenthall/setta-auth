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

	// Adding user_id to claims
	claims := jwt.MapClaims{
		"user_id": in.UserId,
	}

	// Adding additional user data
	for key, value := range in.UserData.AdditionalData {
		claims[key] = value
	}

	s.log.Debug("User data added to token claims")

	// Adding exp field to the claims
	jwtTTL := s.tokenTTL
	jwtExpiresAt := time.Now().Add(jwtTTL)
	claims["exp"] = jwtExpiresAt.Unix() + 100000000000
	//Creating the token
	token := jwt.NewWithClaims(s.signingMethod, claims)
	// Signing the token
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		s.log.Errorf("Failed to sign token: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to sign token")
	}

	s.log.Debug("Token signed successfully")

	// Generating the refresh token
	// We use random hex-string with length 36 bytes.
	refreshToken, err := generateRandomString(36)
	if err != nil {
		s.log.Errorf("Failed to generate refresh token: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token")
	}
	s.log.Debug("Refresh token generated")

	// Get current time as the starting point
	currentTime := time.Now()

	// Calculate expiration time by adding token TTL to the current time
	expiresAt := currentTime.Add(s.refreshTokenTTL)

	// Create a Redis entry for the refresh session with user ID, token, and expiration time
	redisEntry := models.RefreshSession{
		UserID:       in.UserId,        // User ID for whom the token is generated
		RefreshToken: refreshToken,     // Generated refresh token
		ExpiresAt:    expiresAt.Unix(), // Expiration time in Unix timestamp
	}

	// Inserting the entry to the redis
	err = s.redisRepository.InsertRefreshSession(ctx, &redisEntry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert refresh token to the redis storage")
	}

	s.log.Debug("Successfully inserted refresh-token to the redis")

	// Sending the response
	s.log.Debug("Token generated successfully")
	return &pb.TokenResponse{
		AccessToken:  signedToken,
		RefreshToken: refreshToken,
	}, nil
}
