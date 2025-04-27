package analysis

import (
	"web-analyzer/models"
)

var analysisResults = make(map[string]models.AnalysisResult)

func NewAnalysis() *Analysis {

	return &Analysis{}

}

// Storage represents a storage structure.

type Analysis struct {
	analysisData map[string]models.AnalysisResult
}

// StoreAnalysis implements analyzer.Analysis.
func (a Analysis) StoreAnalysis(url string, result models.AnalysisResult) {
	analysisResults[url] = result
}

func (a Analysis) GetAnalysis(url string) (models.AnalysisResult, bool) {

	result, exists := analysisResults[url]
	return result, exists

}
