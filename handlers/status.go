package handlers

import (
	"log/slog"
	"net/http"

	services "web-analyzer/internal/storage"

	"github.com/gin-gonic/gin"
)

func StatusHandler(c *gin.Context) {
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
