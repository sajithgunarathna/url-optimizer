package analyzer

import (
	"web-analyzer/models"
)

type AnalyzerService interface {
	AnalyzePage(url string)
	GetAnalysis(url string) (models.AnalysisResult, bool)
}

type Analysis interface {
	StoreAnalysis(url string, result models.AnalysisResult)
	GetAnalysis(url string) (models.AnalysisResult, bool)
}

type Analyzer struct {
	Storage     Storage
	LinkChecker LinkChecker
	Analysis    Analysis
}

type Storage interface {
	AddSubmittedUrl(url string)
}

type LinkChecker interface {
	IsBroken(url string) bool
}

type AnalysisResult struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}
