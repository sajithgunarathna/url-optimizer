package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"web-analyzer/internal/analyzer"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestDetectHTMLVersion_HTML5(t *testing.T) {
	htmlContent := "<!DOCTYPE html><html><head></head><body></body></html>"
	doc, _ := html.Parse(strings.NewReader(htmlContent))
	version := analyzer.DetectHTMLVersion(doc)
	assert.Equal(t, "HTML5", version)
}

func TestDetectHTMLVersion_Unknown(t *testing.T) {
	htmlContent := "<html><head></head><body></body></html>"
	doc, _ := html.Parse(strings.NewReader(htmlContent))
	version := analyzer.DetectHTMLVersion(doc)
	assert.Equal(t, "No DOCTYPE found", version)
}

func TestIsBrokenLink_FakeURL(t *testing.T) {
	broken := analyzer.IsBrokenLink("http://invalid.url")
	assert.True(t, broken, "Fake URL should be considered broken")
}

func TestIsBrokenLink_ValidURL(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close() // Ensure server shuts down after test

	broken := analyzer.IsBrokenLink(server.URL)
	assert.False(t, broken, "Test server should not be considered broken")
}
