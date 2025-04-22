package test

import (
	"testing"
	services "web-analyzer/internal/storage"
	"web-analyzer/models"

	"github.com/stretchr/testify/assert"
)

func TestStoreAndGetAnalysis(t *testing.T) {
	url := "http://example.com"
	analysis := models.AnalysisResult{Status: "Completed"}

	services.StoreAnalysis(url, analysis)
	result, exists := services.GetAnalysis(url)

	assert.True(t, exists)
	assert.Equal(t, "Completed", result.Status)
}

func TestGetSubmittedUrls(t *testing.T) {
	url := "http://example.com"
	services.AddSubmittedUrl(url)

	urls := services.GetSubmittedUrls()
	assert.Contains(t, urls, url, "Submitted URLs should contain the test URL")
}
