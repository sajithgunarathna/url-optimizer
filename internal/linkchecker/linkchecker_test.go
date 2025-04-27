package linkchecker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDefaultLinkChecker_IsBroken(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockStatusCode int
		mockDelay      time.Duration
		expectedBroken bool
	}{
		{
			name:           "Valid link",
			url:            "http://example.com",
			mockStatusCode: http.StatusOK,
			expectedBroken: false,
		},
		{
			name:           "Broken link with 404",
			url:            "http://example.com/404",
			mockStatusCode: http.StatusNotFound,
			expectedBroken: true,
		},
		{
			name:           "Server error with 500",
			url:            "http://example.com/500",
			mockStatusCode: http.StatusInternalServerError,
			expectedBroken: true,
		},
		{
			name:           "Timeout error",
			url:            "http://example.com/timeout",
			mockDelay:      6 * time.Second,
			expectedBroken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.mockDelay > 0 {
					time.Sleep(tt.mockDelay)
				}
				w.WriteHeader(tt.mockStatusCode)
			}))
			defer server.Close()

			// Replace the URL with the mock server's URL
			url := server.URL
			if tt.url != "http://example.com" {
				url = server.URL + tt.url[len("http://example.com"):]
			}

			// Create a DefaultLinkChecker and test IsBroken
			checker := DefaultLinkChecker{}
			isBroken := checker.IsBroken(url)
			if isBroken != tt.expectedBroken {
				t.Errorf("IsBroken(%q) = %v; want %v", url, isBroken, tt.expectedBroken)
			}
		})
	}
}
