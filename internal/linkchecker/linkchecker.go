package linkchecker

import (
	"log/slog"
	"net/http"
	"time"
)

type LinkChecker interface {
	IsBroken(url string) bool
}

type DefaultLinkChecker struct{}

func (d DefaultLinkChecker) IsBroken(url string) bool {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(url)
	if err != nil || resp.StatusCode >= 400 {
		slog.Debug("Broken or failed link", "url", url, "error", err, "status", resp.StatusCode)
		return true
	}
	return false
}

func NewLinkChecker() LinkChecker {
	return &DefaultLinkChecker{}
}

// Analyzer represents the structure for analyzing HTML documents.
type Analyzer struct {
	LinkChecker LinkChecker
}
