// Package handlers provides HTTP handlers for processing and managing URL analysis requests.
// It includes functionality for analyzing URLs, retrieving submitted URLs, and checking the status of URL analyses.
package handlers

import (
	"log/slog"
	"net/http"
	"regexp"

	"web-analyzer/internal/analyzer"

	"github.com/gin-gonic/gin"
)

// AnalyzeRequest represents the structure of the request body for analyzing a URL.
type AnalyzeRequest struct {
	URL string `json:"url"`
}

// Handler provides HTTP handlers for URL analysis operations.
// It includes methods for analyzing URLs, retrieving submitted URLs, and checking the status of URL analyses.
type Handler struct {
	AnalyzerService analyzer.AnalyzerService
}

// NewHandler creates a new instance of Handler with the provided AnalyzerService.
// It is used to initialize the HTTP handlers for URL analysis operations.
func NewHandler(analyzerService analyzer.AnalyzerService) *Handler {
	return &Handler{
		AnalyzerService: analyzerService,
	}
}

// AnalyzeHandler handles the HTTP request for analyzing a URL.
// It validates the request, checks the URL format, and submits the URL for analysis.
func (h *Handler) AnalyzeHandler(c *gin.Context) {
	var req AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	url := req.URL

	if url == "" {
		slog.Warn("URL parameter is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}
	if !isValidURL(url) {
		slog.Warn("Invalid URL format", "url", url)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	slog.Info("URL submitted for analysis", "url", url)
	go h.AnalyzerService.AnalyzePage(url)

	c.JSON(http.StatusAccepted, gin.H{"message": "URL submitted for analysis"})
}

func isValidURL(url string) bool {
	const urlPattern = `^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(:[0-9]{1,5})?(/.*)?$`
	re := regexp.MustCompile(urlPattern)
	return re.MatchString(url)
}
