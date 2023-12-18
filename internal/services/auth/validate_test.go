package auth

import (
	pb "auth_service/internal/api"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJwtAuthService_ValidateToken(t *testing.T) {
	authService := newTestJwtAuthService(t)

	// Test request with valid token
	req := &pb.ValidateTokenRequest{
		Token: testValidToken,
	}

	// Calling the method under test
	resp, err := authService.ValidateToken(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, true, resp.IsValid)

	// Test request with invalid token
	req = &pb.ValidateTokenRequest{
		Token: testInvalidToken,
	}

	// Calling the method under test
	resp, err = authService.ValidateToken(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, false, resp.IsValid)

}
