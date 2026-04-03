package server

import (
	"context"
	"testing"
)

func TestNewWhiteListMatcher(t *testing.T) {
	whiteList := []string{
		"/health",
		"/api.file.v1.FileService/HealthCheck",
		"/public/",
	}
	matcher := NewWhiteListMatcher(whiteList)

	tests := []struct {
		operation string
		want      bool // true = execute middleware, false = skip
	}{
		{"/health", false},
		{"/api.file.v1.FileService/HealthCheck", false},
		{"/api.file.v1.FileService/UploadFile", true},
		{"/public/anything", false},
		{"/private/anything", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.operation, func(t *testing.T) {
			if got := matcher(context.Background(), tt.operation); got != tt.want {
				t.Errorf("NewWhiteListMatcher() = %v, want %v for operation %v", got, tt.want, tt.operation)
			}
		})
	}
}
