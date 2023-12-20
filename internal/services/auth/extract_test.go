package auth

import (
	pb "auth_service/internal/api"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test extracts data from pre-generated token via authService.ExtractTokenData
func TestJwtAuthService_ExtractTokenData(t *testing.T) {
	authService := newTestJwtAuthService(t)

	req := &pb.ExtractTokenDataRequest{
		Token: testValidToken,
	}

	// Call of ExtractTokenData
	resp, err := authService.ExtractTokenData(context.Background(), req)

	// basic checks
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Data)

	// Checking fields
	assert.Equal(t, "12345", resp.Data["id"])
	assert.Equal(t, "true", resp.Data["booleanParam"])
	assert.Equal(t, "music", resp.Data["word"])

}
