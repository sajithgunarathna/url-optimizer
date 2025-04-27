package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyzer/handlers"
	"web-analyzer/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockStorage for testing
type MockStorage struct{}

// AnalyzePage implements analyzer.AnalyzerService.
func (m *MockStorage) AnalyzePage(url string) {
	panic("unimplemented")
}

// GetAnalysis implements analyzer.AnalyzerService.
func (m *MockStorage) GetAnalysis(url string) (models.AnalysisResult, bool) {
	panic("unimplemented")
}

func (m *MockStorage) GetSubmittedUrls() []string {
	return []string{"http://example.com", "http://another-example.com"}
}

func setupTestHandler() *handlers.Handler {
	mockStorage := &MockStorage{}
	return handlers.NewHandler(mockStorage)
}

func TestUrlsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := setupTestHandler()
	router := gin.Default()
	router.GET("/urls", h.UrlsHandler)

	// Test the /urls endpoint
	req, _ := http.NewRequest(http.MethodGet, "/urls", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and the content of the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http://example.com")
	assert.Contains(t, w.Body.String(), "http://another-example.com")
}
