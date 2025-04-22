package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"web-analyzer/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUrlsHandler_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/urls", handlers.UrlsHandler)

	req, _ := http.NewRequest(http.MethodGet, "/urls", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
