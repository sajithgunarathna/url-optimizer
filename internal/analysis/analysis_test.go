package analysis_test

import (
	"testing"
	"web-analyzer/models"

	"github.com/stretchr/testify/assert"
)

func TestAnalysisResultInitialization(t *testing.T) {
	result := models.AnalysisResult{
		Status:        "Success",
		HTMLVersion:   "HTML5",
		Title:         "Test Page",
		Headings:      map[string]int{"h1": 2, "h2": 3},
		InternalLinks: 5,
		ExternalLinks: 3,
		BrokenLinks:   1,
		LoginForm:     "Present",
		Message:       "Analysis completed successfully",
	}

	assert.Equal(t, "Success", result.Status)
	assert.Equal(t, "HTML5", result.HTMLVersion)
	assert.Equal(t, "Test Page", result.Title)
	assert.Equal(t, 2, result.Headings["h1"])
	assert.Equal(t, 3, result.Headings["h2"])
	assert.Equal(t, 5, result.InternalLinks)
	assert.Equal(t, 3, result.ExternalLinks)
	assert.Equal(t, 1, result.BrokenLinks)
	assert.Equal(t, "Present", result.LoginForm)
	assert.Equal(t, "Analysis completed successfully", result.Message)
}

func TestAnalysisResultEmptyInitialization(t *testing.T) {
	result := models.AnalysisResult{}

	assert.Empty(t, result.Status)
	assert.Empty(t, result.HTMLVersion)
	assert.Empty(t, result.Title)
	assert.Empty(t, result.Headings)
	assert.Equal(t, 0, result.InternalLinks)
	assert.Equal(t, 0, result.ExternalLinks)
	assert.Equal(t, 0, result.BrokenLinks)
	assert.Empty(t, result.LoginForm)
	assert.Empty(t, result.Message)
}

func TestAnalysisResultHeadings(t *testing.T) {
	result := models.AnalysisResult{
		Headings: map[string]int{"h1": 1, "h2": 2, "h3": 0},
	}

	assert.Equal(t, 1, result.Headings["h1"])
	assert.Equal(t, 2, result.Headings["h2"])
	assert.Equal(t, 0, result.Headings["h3"])
}
