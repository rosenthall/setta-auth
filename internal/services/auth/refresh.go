package auth

import (
	pb "auth_service/internal/api"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *JwtAuthService) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	s.log.Debug("Refreshing token.")

	// Basic validation
	if in.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "RefreshToken must contain a refresh token.")
	}

	// Getting refresh session by refresh token
	refreshSession, err := s.redisRepository.GetRefreshSession(ctx, in.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Cannot find refresh session for the given token.")
	}

	// Checking if refresh token expired
	if time.Now().Unix() > refreshSession.ExpiresAt {
		return nil, status.Error(codes.Unauthenticated, "Refresh token is expired.")
	}

	// Comparing the provided refresh token with the one stored in the session
	if refreshSession.RefreshToken != in.RefreshToken {
		return nil, status.Error(codes.Unauthenticated, "Invalid refresh session.")
	}

	// Deleting old refresh session entry
	err = s.redisRepository.DeleteRefreshSession(ctx, in.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot delete old refresh session: %s", err)
	}

	// Generate new access and refresh tokens
	newAccessToken, newRefreshToken, err := s.generateTokens(ctx, refreshSession.UserID, nil)
	if err != nil {
		return nil, err
	}

	s.log.Debug("New tokens generated successfully")

	// Return new tokens to the client
	return &pb.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
