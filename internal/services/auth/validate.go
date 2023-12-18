package auth

import (
	pb "auth_service/internal/api"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *JwtAuthService) ValidateToken(_ context.Context, in *pb.ValidateTokenRequest) (*pb.TokenValidationResponse, error) {
	s.log.Debug("Validating token")
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "Token must contain a JWT token")
	}

	token, err := jwt.Parse(in.Token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "unexpected signing method: %v", token.Header["alg"])
		}
		// Return the public key for signature verification
		return s.publicKey, nil
	})

	if err != nil {
		// In case of error, consider the token invalid
		s.log.Errorf("Error while parsing token for validating : %s", err)
		return &pb.TokenValidationResponse{IsValid: false}, nil
	}

	// Check if the token is valid and not expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check for expiration
		if exp, ok := claims["exp"].(float64); ok {
			// Convert to time.Time and compare with the current time
			if time.Unix(int64(exp), 0).After(time.Now()) {
				// Token is valid and not expired
				s.log.Debug("Token is valid")
				return &pb.TokenValidationResponse{IsValid: true}, nil
			}
		}
	}

	// Token is either invalid or expired
	s.log.Debug("Token is invalid")
	return &pb.TokenValidationResponse{IsValid: false}, nil
}
