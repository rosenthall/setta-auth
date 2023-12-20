package auth

import (
	pb "auth_service/internal/api"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// This token will never expire due of huge `exp` value
	testValidToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJib29sZWFuUGFyYW0iOiJ0cnVlIiwiZXhwIjoxMDE3MDI4OTYwNjYsImlkIjoiMTIzNDUiLCJ1c2VyX2lkIjoid2kzOCIsIndvcmQiOiJtdXNpYyJ9.V6hZE_-2nrM0A7KhT5TxmI_xMtll8gZnqjqWEOSgRYXUVWTzf8xpVIN1WU4-hgorEM_eymxfJSyAEN_5x4VjQEzxXJ8vaelrGFr6CBoG0h-n6V96qw8ugMlNTvCRmnoVSh5l5ozuWl3UHxaxWuxSyrSbHQeNFNLb5qK-NY5DHflWmsnHHZ75KekxFBRYMmOdqIcSrQRMF_98aJT3yeH9zIdr9Xg0Bjnc9c7xI0HFmJtpsZnDSUR3XNqpJIEyPx59l4s9ANALmHynmQIZ-48CPmPPCBjGu28zBN-Z8D6gwj8AatMfz9BBZFdQTQuPZ7PEdfrhBfMEFEth7-TsvajnQA"
	// This token is just invalid
	testInvalidToken = "eyJhbGciOiJSUzI1NiIsInR5CI6pXVJ.eleHAiOjEwMTcwMjg5NzY2MCwidXNlcl9pZCI6IndMgi.bWj2lzXuydPQjNreD-Kt5bm3hVScm9pdiI4jGj08N-nh9ey7_ONxnJkfdWnPTcOyiqM3DIHp0qajZxdzcr6G558eAwqoD77opjtSSiZLYW37wtDpK2Z_l3kUoMqzSuESTFK7bVz42Rz7vUoUqp8Y6YemcmonbkbY4BFSI72cwJKrpUgg5vOUvpSYEthQ_cTzZWWWCkrSuDcpTszS22TI6cPAapln5lHNWhwG7n6vnPPZRNtWLL2ZTh9ohubiy4pMjcSgAe_ANThKbsUkrYw7OEKSkudwZwx6975hhVbNtbZ0XRXONrVMQbm_Z4JRvXmKui1QTu0gF"
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
