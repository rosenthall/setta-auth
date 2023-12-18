package auth

import (
	pb "auth_service/internal/api"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testValidToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJib29sZWFuUGFyYW0iOiJ0cnVlIiwiZXhwIjoxNzAyODM4Mjc4LCJpZCI6IjEyMzQ1IiwidXNlcl9pZCI6IndpMzgiLCJ3b3JkIjoibXVzaWMifQ.rDxHY5_2YHi0EazRdJ-YRjTlEB7aCn6vYbPUKhnMfMRx43UZISTZ6Ei1-i0nJgEBSmE-d9Q5-VBXuqB90iGyG8pkMeWrlDXY_z7yawEe07Zn8ON2lnt2OyW7OylzS8tyArqArNFU5d90boAWE1hIbNX2Y1LlXBV6nayOvW_SPnc-UG4q-wv78v6tcshSSOEXfKwOfkO0eL9GPDX6v-X3EogWCHHwkWgq1LgOTU3mHk4EQ3aYv2XMTA-m8CCGhlO5vVXVsRUlMRwVgiSAjkEmlqhJ2vMUSWn9rz5MkU1gOKJhubrc6cZyMXcS0TW6_95Dfv5v3uUtewIZXlzLFXz7BA"
)

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
