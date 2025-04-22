package handlers

import (
	"log/slog"
	"net/http"
	"regexp"

	"web-analyzer/internal/analyzer"

	services "web-analyzer/internal/storage"

	"github.com/gin-gonic/gin"
)

type AnalyzeRequest struct {
	URL string `json:"url"`
}

type Handler struct {
	Analyzer analyzer.AnalyzerService
}

func NewHandler(analyzer analyzer.AnalyzerService) *Handler {
	return &Handler{
		Analyzer: analyzer,
	}
}

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
	go h.Analyzer.AnalyzePage(url)

	c.JSON(http.StatusAccepted, gin.H{"message": "URL submitted for analysis"})
}

func isValidURL(url string) bool {
	const urlPattern = `^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(:[0-9]{1,5})?(/.*)?$`
	re := regexp.MustCompile(urlPattern)
	return re.MatchString(url)
}

func (h *Handler) UrlsHandler(c *gin.Context) {
	// Example stub logic
	c.JSON(http.StatusOK, gin.H{"message": "List of URLs"})
}

func (h *Handler) StatusHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		slog.Warn("Missing URL parameter in status check")
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	analysis, exists := services.GetAnalysis(url)
	if !exists {
		slog.Info("Analysis not found", "url", url)
		c.JSON(http.StatusNotFound, gin.H{"error": "Analysis not found"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
