package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyzer/handlers"
	"web-analyzer/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockAnalyzerService implements analyzer.AnalyzerService
type mockAnalyzerService struct{}

func (m *mockAnalyzerService) AnalyzePage(url string) {
	// Mock implementation of AnalyzePage
}

func (m *mockAnalyzerService) GetAnalysis(url string) (models.AnalysisResult, bool) {
	// Mock implementation of GetAnalysis
	return models.AnalysisResult{
		Title: "Mock Title",
		// Description field removed as it does not exist in models.AnalysisResult
		// Keywords field removed as it does not exist in models.AnalysisResult
	}, true
}

func TestAnalyzeHandler_ValidURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyzerService{}
	handler := handlers.NewHandler(mockService)

	router := gin.Default()
	router.POST("/analyze", handler.AnalyzeHandler)

	reqBody := handlers.AnalyzeRequest{URL: "https://example.com"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "URL submitted for analysis", response["message"])
}

func TestAnalyzeHandler_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyzerService{}
	handler := handlers.NewHandler(mockService)

	router := gin.Default()
	router.POST("/analyze", handler.AnalyzeHandler)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Invalid request", response["error"])
}

func TestAnalyzeHandler_MissingURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyzerService{}
	handler := handlers.NewHandler(mockService)

	router := gin.Default()
	router.POST("/analyze", handler.AnalyzeHandler)

	reqBody := handlers.AnalyzeRequest{URL: ""}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "URL parameter is required", response["error"])
}

func TestAnalyzeHandler_InvalidURLFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyzerService{}
	handler := handlers.NewHandler(mockService)

	router := gin.Default()
	router.POST("/analyze", handler.AnalyzeHandler)

	reqBody := handlers.AnalyzeRequest{URL: "invalid-url"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Invalid URL", response["error"])
}
