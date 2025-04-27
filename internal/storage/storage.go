package services

import (
	"log/slog"
)

var submittedUrls = make(map[string]struct{})

func GetSubmittedUrls() []string {
	urls := make([]string, 0, len(submittedUrls))
	for url := range submittedUrls {
		urls = append(urls, url)
	}
	return urls
}

// NewStorage creates and returns a new instance of Storage.

func NewStorage() *Storage {

	return &Storage{}

}

// Storage represents a storage structure.

type Storage struct {
	submittedUrls map[string]bool
}

// AddSubmittedUrl implements analyzer.Storage.
func (s Storage) AddSubmittedUrl(url string) {
	submittedUrls[url] = struct{}{}
	slog.Info("New URL added to submissions", "url", url)
}

type LinkChecker interface {
	CheckLink(url string) bool
}

type Analyzer struct {
	Storage     Storage
	LinkChecker LinkChecker
}

func (a *Analyzer) GetSubmittedUrls() []string {
	urls := make([]string, 0, len(submittedUrls))
	for url := range submittedUrls {
		urls = append(urls, url)
	}
	return urls
}
