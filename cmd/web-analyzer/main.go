package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	"web-analyzer/handlers"
	"web-analyzer/internal/analysis"
	"web-analyzer/internal/analyzer"
	"web-analyzer/internal/server"

	"web-analyzer/internal/linkchecker"
	services "web-analyzer/internal/storage"

	"github.com/gin-contrib/pprof"
	"github.com/rs/cors"
)

func main() {
	// Create instances of Storage and LinkChecker
	storage := services.NewStorage()
	linkChecker := linkchecker.NewLinkChecker()
	analysis := analysis.NewAnalysis() // Pass storage to analysis

	// Pass the required arguments to NewAnalyzerService
	analyzerService := analyzer.NewAnalyzerService(*storage, linkChecker, *analysis)
	h := handlers.NewHandler(analyzerService)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize handlers
	// h is already initialized with NewHandler

	// Pass the handler to SetupRouter
	r := server.SetupRouter(h)

	pprof.Register(r)

	// Enable CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust based on frontend
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	finalHandler := corsMiddleware.Handler(r)

	logger.Info("Server started on :8080")
	err := http.ListenAndServe(":8080", finalHandler)
	if err != nil {
		logger.Error("Server failed to start", "error", err)
	}
}
