package models

type AnalysisResult struct {
	Status        string         `json:"Status"`
	HTMLVersion   string         `json:"HTML Version"`
	Title         string         `json:"Title"`
	Headings      map[string]int `json:"Headings"`
	InternalLinks int            `json:"Internal Links"`
	ExternalLinks int            `json:"External Links"`
	BrokenLinks   int            `json:"Broken Links"`
	LoginForm     string         `json:"Login Form"`
	Message       string         `json:"Message,omitempty"`
}
