package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"web-analyzer/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStatusHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/status", handlers.StatusHandler)

	req, _ := http.NewRequest(http.MethodGet, "/status?url=http://notfound.com", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
