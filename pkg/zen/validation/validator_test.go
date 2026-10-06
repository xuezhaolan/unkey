package validation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaIsValid(t *testing.T) {
	_, err := New()
	require.NoError(t, err)
}

func TestRatelimitLimitDurationMaximumUsesDecimalNotation(t *testing.T) {
	validator, err := New()
	require.NoError(t, err)

	body := strings.NewReader(`{
		"namespace": "api.requests",
		"identifier": "user_123",
		"limit": 100,
		"duration": 2592000001
	}`)
	req := httptest.NewRequest(http.MethodPost, "/v2/ratelimit.limit", body)
	req.Header.Set("Authorization", "Bearer test")
	req.Header.Set("Content-Type", "application/json")

	res, valid := validator.Validate(context.Background(), req)
	require.False(t, valid)

	var durationError string
	for _, validationError := range res.Error.Errors {
		if validationError.Location == "/properties/duration/maximum" {
			durationError = validationError.Message
			break
		}
	}

	require.Equal(t, "maximum: got 2592000001, want 2592000000", durationError)
	require.NotContains(t, durationError, "e+")
}
