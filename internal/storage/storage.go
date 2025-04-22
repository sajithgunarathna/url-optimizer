package services

import (
	"log/slog"
	"web-analyzer/models"
)

var analysisResults = make(map[string]models.AnalysisResult)
var submittedUrls = make(map[string]struct{})

func StoreAnalysis(url string, result models.AnalysisResult) {
	analysisResults[url] = result
}

func GetAnalysis(url string) (models.AnalysisResult, bool) {
	result, exists := analysisResults[url]
	return result, exists
}

func GetSubmittedUrls() []string {
	urls := make([]string, 0, len(submittedUrls))
	for url := range submittedUrls {
		urls = append(urls, url)
	}
	return urls
}

func AddSubmittedUrl(url string) {
	submittedUrls[url] = struct{}{}
	slog.Info("New URL added to submissions", "url", url)
}
