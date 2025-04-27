package analyzer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyzer/models"
)

type mockStorage struct {
	submittedUrls map[string]bool
}

func (m *mockStorage) AddSubmittedUrl(url string) {
	m.submittedUrls[url] = true
}

type mockLinkChecker struct {
	brokenLinks map[string]bool
}

func (m *mockLinkChecker) IsBroken(url string) bool {
	return m.brokenLinks[url]
}

type mockAnalysis struct {
	analysisResults map[string]models.AnalysisResult
}

func (m *mockAnalysis) GetAnalysis(url string) (models.AnalysisResult, bool) {
	result, exists := m.analysisResults[url]
	return result, exists
}

func (m *mockAnalysis) StoreAnalysis(url string, result models.AnalysisResult) {
	m.analysisResults[url] = result
}

func TestAnalyzePage_Success(t *testing.T) {
	mockStorage := &mockStorage{submittedUrls: make(map[string]bool)}
	mockLinkChecker := &mockLinkChecker{brokenLinks: make(map[string]bool)}
	mockAnalysis := &mockAnalysis{analysisResults: make(map[string]models.AnalysisResult)}

	service := DefaultAnalyzerService{
		Analyzer: &Analyzer{
			Storage:     mockStorage,
			LinkChecker: mockLinkChecker,
			Analysis:    mockAnalysis,
		},
	}

	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<h1>Header 1</h1>
			<h2>Header 2</h2>
			<a href="/internal">Internal Link</a>
			<a href="http://external.com">External Link</a>
			<form>
				<input type="password" />
			</form>
		</body>
		</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	service.AnalyzePage(server.URL)

	result, exists := mockAnalysis.GetAnalysis(server.URL)
	if !exists {
		t.Fatalf("Expected analysis result to exist for URL: %s", server.URL)
	}

	if result.Title != "Test Page" {
		t.Errorf("Expected title to be 'Test Page', got '%s'", result.Title)
	}

	if result.Headings["h1"] != 1 || result.Headings["h2"] != 1 {
		t.Errorf("Expected headings count to be correct, got %v", result.Headings)
	}

	if result.InternalLinks != 1 || result.ExternalLinks != 1 {
		t.Errorf("Expected internal and external links count to be 1, got %d and %d", result.InternalLinks, result.ExternalLinks)
	}

	if result.LoginForm != "Present" {
		t.Errorf("Expected login form to be 'Present', got '%s'", result.LoginForm)
	}
}

func TestAnalyzePage_FetchError(t *testing.T) {
	mockStorage := &mockStorage{submittedUrls: make(map[string]bool)}
	mockLinkChecker := &mockLinkChecker{brokenLinks: make(map[string]bool)}
	mockAnalysis := &mockAnalysis{analysisResults: make(map[string]models.AnalysisResult)}

	service := DefaultAnalyzerService{
		Analyzer: &Analyzer{
			Storage:     mockStorage,
			LinkChecker: mockLinkChecker,
			Analysis:    mockAnalysis,
		},
	}

	invalidURL := "http://invalid-url"

	service.AnalyzePage(invalidURL)

	_, exists := mockAnalysis.GetAnalysis(invalidURL)
	if exists {
		t.Fatalf("Expected no analysis result for invalid URL: %s", invalidURL)
	}
}

func TestAnalyzePage_HTMLParseError(t *testing.T) {
	mockStorage := &mockStorage{submittedUrls: make(map[string]bool)}
	mockLinkChecker := &mockLinkChecker{brokenLinks: make(map[string]bool)}
	mockAnalysis := &mockAnalysis{analysisResults: make(map[string]models.AnalysisResult)}

	service := DefaultAnalyzerService{
		Analyzer: &Analyzer{
			Storage:     mockStorage,
			LinkChecker: mockLinkChecker,
			Analysis:    mockAnalysis,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invalid HTML"))
	}))
	defer server.Close()

	service.AnalyzePage(server.URL)

	_, exists := mockAnalysis.GetAnalysis(server.URL)
	if exists {
		t.Fatalf("Expected no analysis result for invalid HTML")
	}
}
