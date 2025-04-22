package server

import (
	"log/slog"
	"net/http"
	"web-analyzer/handlers"

	"github.com/gin-gonic/gin"
)

func Start(h *handlers.Handler) error {
	r := SetupRouter(h)
	return r.Run(":8080")
}

func SetupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()
	setupPprof(r)

	r.POST("/analyze", h.AnalyzeHandler)
	r.GET("/status", h.StatusHandler)
	r.GET("/urls", h.UrlsHandler)

	return r
}

func setupPprof(router *gin.Engine) {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			slog.Error("pprof server failed", "error", err)
		}
	}()
}
