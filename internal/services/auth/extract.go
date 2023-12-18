package auth

import (
	pb "auth_service/internal/api"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *JwtAuthService) ExtractTokenData(_ context.Context, in *pb.ExtractTokenDataRequest) (*pb.TokenDataResponse, error) {
	// Check if the token is present
	if in.Token == "" {
		s.log.Warn("Token is empty")
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Parse the token
	token, err := jwt.Parse(in.Token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			s.log.Warnf("Unexpected signing method: %v", token.Header["alg"])
			return nil, status.Errorf(codes.Unauthenticated, "unexpected signing method: %v", token.Header["alg"])
		}

		// Return the public key for signature verification
		return s.publicKey, nil
	})

	if err != nil {
		s.log.Errorf("Invalid token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Validate the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		s.log.Debug("Token is valid, extracting data")

		// Extract data from the token
		data := make(map[string]string)
		for key, value := range claims {
			// Convert all values to strings
			data[key] = fmt.Sprintf("%v", value)
			s.log.Debug(fmt.Sprintf("%v", value))
		}

		s.log.Debugf("Extracted data: %v", data)

		// Form the response
		return &pb.TokenDataResponse{
			Data: data,
		}, nil
	} else {
		s.log.Warn("Invalid token claims")
		return nil, status.Error(codes.Unauthenticated, "invalid token claims")
	}
}
