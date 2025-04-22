package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"web-analyzer/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockAnalyzerService struct{}

func (m *MockAnalyzerService) AnalyzePage(url string) {
	// Mock implementation
}

func TestAnalyzeHandler_ValidURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnalyzer := &MockAnalyzerService{}
	handler := handlers.NewHandler(mockAnalyzer)

	router := gin.Default()
	router.POST("/analyze", handler.AnalyzeHandler)

	reqBody := handlers.AnalyzeRequest{URL: "http://example.com"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.JSONEq(t, `{"message": "URL submitted for analysis"}`, rec.Body.String())
}

func TestAnalyzeHandler_InvalidURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnalyzer := &MockAnalyzerService{}
	handler := handlers.NewHandler(mockAnalyzer)

	router := gin.Default()
	router.POST("/analyze", handler.AnalyzeHandler)

	reqBody := handlers.AnalyzeRequest{URL: "invalid-url"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error": "Invalid URL"}`, rec.Body.String())
}

func TestAnalyzeHandler_MissingURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnalyzer := &MockAnalyzerService{}
	handler := handlers.NewHandler(mockAnalyzer)

	router := gin.Default()
	router.POST("/analyze", handler.AnalyzeHandler)

	reqBody := handlers.AnalyzeRequest{URL: ""}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error": "URL parameter is required"}`, rec.Body.String())
}

func TestUrlsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAnalyzer := &MockAnalyzerService{}
	handler := handlers.NewHandler(mockAnalyzer)

	router := gin.Default()
	router.GET("/urls", handler.UrlsHandler)

	req, _ := http.NewRequest(http.MethodGet, "/urls", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message": "List of URLs"}`, rec.Body.String())
}
