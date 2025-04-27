package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyzer/handlers"
	"web-analyzer/internal/server"

	"web-analyzer/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockAnalyzerService for testing
type MockAnalyzerService struct{}

func (m *MockAnalyzerService) AnalyzePage(url string) {}

func (m *MockAnalyzerService) GetAnalysis(url string) (models.AnalysisResult, bool) {
	return models.AnalysisResult{}, false
}

func setupTestHandler() *handlers.Handler {
	service := &MockAnalyzerService{}
	return handlers.NewHandler(service)
}

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := setupTestHandler()
	router := server.SetupRouter(h)

	// Test /analyze endpoint
	req1, _ := http.NewRequest(http.MethodPost, "/analyze", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Contains(t, []int{http.StatusBadRequest, http.StatusAccepted}, w1.Code)

	// Test /status endpoint
	req2, _ := http.NewRequest(http.MethodGet, "/status?url=http://example.com", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	// Test /urls endpoint
	req3, _ := http.NewRequest(http.MethodGet, "/urls", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
