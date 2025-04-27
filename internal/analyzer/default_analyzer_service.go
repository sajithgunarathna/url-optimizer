package analyzer

import (
	"log/slog"
	"net/http"
	"strings"
	"web-analyzer/models"

	"golang.org/x/net/html"
)

const inProgress = "In progress"

type DefaultAnalyzerService struct {
	Analyzer *Analyzer
}

func (d DefaultAnalyzerService) GetAnalysis(url string) (models.AnalysisResult, bool) {
	return d.Analyzer.Analysis.GetAnalysis(url)
}

func (d DefaultAnalyzerService) AnalyzePage(url string) {
	d.Analyzer.Storage.AddSubmittedUrl(url)
	slog.Info("AnalyzePage called", "url", url)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		slog.Error("Error fetching or invalid status", "url", url, "status", resp.StatusCode, "error", err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		slog.Error("Error parsing HTML", "error", err)
		return
	}

	inProgressResult := models.AnalysisResult{
		Status:    inProgress,
		Headings:  make(map[string]int),
		LoginForm: "Not Present",
	}

	d.Analyzer.Analysis.StoreAnalysis(url, inProgressResult)

	result := d.Analyzer.AnalyzeHTML(doc, url)
	slog.Info("Analysis Result", "result", result)
	d.Analyzer.Analysis.StoreAnalysis(url, result)
}

func NewAnalyzerService(storage Storage, linkChecker LinkChecker, analysis Analysis) AnalyzerService {
	return &DefaultAnalyzerService{
		Analyzer: &Analyzer{
			Storage:     storage,
			LinkChecker: linkChecker,
			Analysis:    analysis,
		},
	}
}

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
