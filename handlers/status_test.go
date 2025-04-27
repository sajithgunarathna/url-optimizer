package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyzer/handlers"
	"web-analyzer/internal/analyzer"
	"web-analyzer/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockAnalyzerService implements the AnalyzerService interface for testing.
type MockAnalyzerService struct {
	analysisData map[string]models.AnalysisResult
}

func (m *MockAnalyzerService) AnalyzePage(url string) {}
func (m *MockAnalyzerService) GetAnalysis(url string) (models.AnalysisResult, bool) {
	result, ok := m.analysisData[url]
	return result, ok
}

func setupRouter(service analyzer.AnalyzerService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	h := handlers.NewHandler(service)
	router.GET("/status", h.StatusHandler)
	return router
}

func TestStatusHandler_MissingURL(t *testing.T) {
	service := &MockAnalyzerService{}
	router := setupRouter(service)

	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestStatusHandler_AnalysisNotFound(t *testing.T) {
	service := &MockAnalyzerService{
		analysisData: make(map[string]models.AnalysisResult),
	}
	router := setupRouter(service)

	req, _ := http.NewRequest(http.MethodGet, "/status?url=http://example.com", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestStatusHandler_AnalysisFound(t *testing.T) {
	service := &MockAnalyzerService{
		analysisData: map[string]models.AnalysisResult{
			"http://example.com": {
				Status: "Completed",
			},
		},
	}
	router := setupRouter(service)

	req, _ := http.NewRequest(http.MethodGet, "/status?url=http://example.com", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Completed")
}
