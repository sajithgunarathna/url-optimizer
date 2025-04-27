package handlers

import (
	"log/slog"
	"net/http"

	"web-analyzer/internal/analyzer"

	"github.com/gin-gonic/gin"
)

type DefaultAnalyzerService struct {
	Analyzer *analyzer.Analyzer
}

// StatusHandler handles the HTTP request for checking the status of a URL analysis.
// It retrieves the analysis result for the specified URL and returns it in the response.
func (h *Handler) StatusHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		slog.Warn("Missing URL parameter in status check")
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	analysis, exists := h.AnalyzerService.GetAnalysis(url)
	if !exists {
		slog.Info("Analysis not found", "url", url)
		c.JSON(http.StatusNotFound, gin.H{"error": "Analysis not found"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
