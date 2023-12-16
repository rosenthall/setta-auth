package auth

import (
	pb "auth_service/internal/api"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJwtAuthService_GenerateToken(t *testing.T) {
	authService := newTestJwtAuthService(t)

	// getting mock object from redisRepository
	mockRedisRepo := authService.redisRepository.(*MockRefreshSessionsRepository)

	// Test request
	req := &pb.GenerateTokenRequest{
		UserId: "test_user",
		UserData: &pb.UserData{
			AdditionalData: map[string]string{
				"role": "user",
			},
		},
	}

	// Calling the method under test
	resp, err := authService.GenerateToken(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	// Asserting expectations on the mock
	mockRedisRepo.AssertExpectations(t)
}
