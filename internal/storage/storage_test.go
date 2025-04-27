package services

import (
	"testing"
)

func TestAddSubmittedUrl(t *testing.T) {
	storage := NewStorage()

	url := "http://example.com"
	storage.AddSubmittedUrl(url)

	submittedUrls := GetSubmittedUrls()
	if len(submittedUrls) != 1 {
		t.Errorf("expected 1 URL, got %d", len(submittedUrls))
	}

	if submittedUrls[0] != url {
		t.Errorf("expected URL %s, got %s", url, submittedUrls[0])
	}
}

func TestGetSubmittedUrls(t *testing.T) {
	storage := NewStorage()

	url1 := "http://example1.com"
	url2 := "http://example2.com"

	storage.AddSubmittedUrl(url1)
	storage.AddSubmittedUrl(url2)

	submittedUrls := GetSubmittedUrls()
	if len(submittedUrls) != 2 {
		t.Errorf("expected 2 URLs, got %d", len(submittedUrls))
	}

	expectedUrls := map[string]bool{
		url1: true,
		url2: true,
	}

	for _, url := range submittedUrls {
		if !expectedUrls[url] {
			t.Errorf("unexpected URL found: %s", url)
		}
	}
}

func TestAnalyzer_GetSubmittedUrls(t *testing.T) {
	storage := NewStorage()
	analyzer := Analyzer{Storage: *storage}

	url1 := "http://example1.com"
	url2 := "http://example2.com"

	storage.AddSubmittedUrl(url1)
	storage.AddSubmittedUrl(url2)

	submittedUrls := analyzer.GetSubmittedUrls()
	if len(submittedUrls) != 2 {
		t.Errorf("expected 2 URLs, got %d", len(submittedUrls))
	}

	expectedUrls := map[string]bool{
		url1: true,
		url2: true,
	}

	for _, url := range submittedUrls {
		if !expectedUrls[url] {
			t.Errorf("unexpected URL found: %s", url)
		}
	}
}
