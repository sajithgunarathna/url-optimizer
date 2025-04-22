package analyzer

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"web-analyzer/models"

	services "web-analyzer/internal/storage"

	"golang.org/x/net/html"
)

// AnalyzerService is the interface used by HTTP handlers.
type AnalyzerService interface {
	AnalyzePage(url string)
}

type DefaultAnalyzerService struct {
	Analyzer *Analyzer
}

func (d DefaultAnalyzerService) AnalyzePage(url string) {

	slog.Info("AnalyzePage called", "url", url)

	// Step 1: Fetch the HTML content from the URL
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Failed to fetch URL, status code: %d\n", resp.StatusCode)
		return
	}

	// Step 2: Parse the HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		slog.Error("Error parsing HTML: %v\n", err)
		return
	}

	analyseInProgress := models.AnalysisResult{
		Status:        inProgress,
		HTMLVersion:   "",
		Title:         "",
		Headings:      make(map[string]int),
		InternalLinks: 0,
		ExternalLinks: 0,
		BrokenLinks:   0,
		LoginForm:     "Not Present",
	}

	services.StoreAnalysis(url, analyseInProgress)
	// Step 3: Analyze the parsed HTML using the Analyzer
	analysis := d.Analyzer.AnalyzeHTML(doc, url)
	slog.Info("Analysis Result", "result", analysis)
	services.StoreAnalysis(url, analysis)

}
func NewAnalyzerService(storage Storage, linkChecker LinkChecker) AnalyzerService {
	return &DefaultAnalyzerService{
		Analyzer: &Analyzer{
			Storage:     storage,
			LinkChecker: linkChecker,
		},
	}
}

// ---------------- Interfaces ---------------- //

type StorageInterface interface {
	AddSubmittedUrl(url string)
	GetAnalysis(url string) (models.AnalysisResult, bool)
	StoreAnalysis(url string, result models.AnalysisResult)
}

type LinkChecker interface {
	IsBroken(url string) bool
}

// ---------------- Implementation ---------------- //

type Analyzer struct {
	Storage     Storage
	LinkChecker LinkChecker
}

const inProgress string = "In progress"

func (a *Analyzer) AnalyzePage(url string) {
	a.Storage.AddSubmittedUrl(url)

	if analysisResults, exists := a.Storage.GetAnalysis(url); exists {
		if analysisResults.Status == inProgress {
			slog.Info("Analysis already in progress", "url", url)
			return
		}
	}

	slog.Info("Starting analysis", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Failed to fetch URL", "url", url, "error", err)
		a.Storage.StoreAnalysis(url, models.AnalysisResult{Status: "Error", Message: "Failed to fetch URL"})
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		slog.Error("Failed to parse HTML", "url", url, "error", err)
		a.Storage.StoreAnalysis(url, models.AnalysisResult{Status: "Error", Message: "Failed to parse HTML"})
		return
	}

	a.Storage.StoreAnalysis(url, models.AnalysisResult{
		Status:        inProgress,
		Headings:      make(map[string]int),
		InternalLinks: 0,
		ExternalLinks: 0,
		BrokenLinks:   0,
		LoginForm:     "Not Present",
	})

	analysis := a.AnalyzeHTML(doc, url)
	a.Storage.StoreAnalysis(url, analysis)
	slog.Info("Analysis completed", "url", url, "status", analysis.Status)
}

func IsBrokenLink(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(url)
	if err != nil {
		slog.Debug("Link check failed", "url", url, "error", err)
		return true
	}
	if resp.StatusCode >= 400 {
		slog.Debug("Broken link detected", "url", url, "status", resp.StatusCode)
		return true
	}
	return false
}

// ---------------- HTML Analyzer ---------------- //

func (a *Analyzer) AnalyzeHTML(doc *html.Node, baseURL string) models.AnalysisResult {
	var title string
	headings := map[string]int{}
	internalLinks, externalLinks, brokenLinks := 0, 0, 0
	loginForm := "Not Present"

	seenLinks := make(map[string]bool)

	var hasPasswordInput func(*html.Node) bool
	hasPasswordInput = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "password" {
					return true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if hasPasswordInput(c) {
				return true
			}
		}
		return false
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				headings[n.Data]++
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href := attr.Val
						fullURL := href
						isExternal := strings.HasPrefix(href, "http")

						if !isExternal {
							internalLinks++
							fullURL = baseURL + href
						} else {
							externalLinks++
						}

						if !seenLinks[fullURL] {
							seenLinks[fullURL] = true
							if a.LinkChecker.IsBroken(fullURL) {
								brokenLinks++
							}
						}
					}
				}
			case "form":
				if hasPasswordInput(n) {
					loginForm = "Present"
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	return models.AnalysisResult{
		Status:        "Completed",
		HTMLVersion:   DetectHTMLVersion(doc),
		Title:         title,
		Headings:      headings,
		InternalLinks: internalLinks,
		ExternalLinks: externalLinks,
		BrokenLinks:   brokenLinks,
		LoginForm:     loginForm,
	}
}

// ---------------- HTML Version Detection ---------------- //

func DetectHTMLVersion(doc *html.Node) string {
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.DoctypeNode {
			if strings.ToLower(c.Data) != "html" {
				return "Unknown document type"
			}
			if len(c.Attr) == 0 {
				return "HTML5"
			}
			for _, attr := range c.Attr {
				if attr.Key == "public" {
					if version, ok := htmlVersions[strings.ToUpper(attr.Val)]; ok {
						return version
					}
					return "Unknown HTML version"
				}
			}
			return "Unknown HTML version"
		}
	}
	return "No DOCTYPE found"
}

var htmlVersions = map[string]string{
	"-//W3C//DTD HTML 4.01//EN":              "HTML 4.01 Strict",
	"-//W3C//DTD HTML 4.01 TRANSITIONAL//EN": "HTML 4.01 Transitional",
	"-//W3C//DTD HTML 4.01 FRAMESET//EN":     "HTML 4.01 Frameset",
	"-//W3C//DTD HTML 4.0//EN":               "HTML 4.0 Strict",
	"-//W3C//DTD HTML 4.0 TRANSITIONAL//EN":  "HTML 4.0 Transitional",
	"-//W3C//DTD HTML 4.0 FRAMESET//EN":      "HTML 4.0 Frameset",
	"-//W3C//DTD HTML 3.2 FINAL//EN":         "HTML 3.2",
	"-//IETF//DTD HTML//EN":                  "HTML 2.0",
	"-//W3C//DTD XHTML 1.0 STRICT//EN":       "XHTML 1.0 Strict",
	"-//W3C//DTD XHTML 1.0 TRANSITIONAL//EN": "XHTML 1.0 Transitional",
	"-//W3C//DTD XHTML 1.0 FRAMESET//EN":     "XHTML 1.0 Frameset",
	"-//W3C//DTD XHTML 1.1//EN":              "XHTML 1.1",
}

// Storage represents a simple storage structure.

type Storage struct {
	submittedUrls map[string]bool
	analysisData  map[string]models.AnalysisResult
}

// StoreAnalysis stores the analysis result for a given URL.
func (s *Storage) StoreAnalysis(url string, result models.AnalysisResult) {
	if s.analysisData == nil {
		s.analysisData = make(map[string]models.AnalysisResult)
	}
	s.analysisData[url] = result
}

// GetAnalysis retrieves the analysis result for a given URL.
func (s *Storage) GetAnalysis(url string) (models.AnalysisResult, bool) {
	if s.analysisData == nil {
		s.analysisData = make(map[string]models.AnalysisResult)
	}
	result, exists := s.analysisData[url]
	return result, exists
}

// AddSubmittedUrl adds a URL to the submitted URLs map.
func (s *Storage) AddSubmittedUrl(url string) {
	if s.submittedUrls == nil {
		s.submittedUrls = make(map[string]bool)
	}
	s.submittedUrls[url] = true
}

// NewStorage creates and returns a new instance of Storage.

func NewStorage() *Storage {

	return &Storage{}

}

// NewLinkChecker creates and returns a new instance of LinkChecker.

type DefaultLinkChecker struct{}

func (d DefaultLinkChecker) IsBroken(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(url)
	if err != nil {
		slog.Debug("Link check failed", "url", url, "error", err)
		return true
	}
	if resp.StatusCode >= 400 {
		slog.Debug("Broken link detected", "url", url, "status", resp.StatusCode)
		return true
	}
	return false
}

func NewLinkChecker() LinkChecker {
	return &DefaultLinkChecker{}
}
