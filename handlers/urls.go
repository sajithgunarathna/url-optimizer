package handlers

import (
	"net/http"
	services "web-analyzer/internal/storage"

	"github.com/gin-gonic/gin"
)

// UrlsHandler handles the HTTP request for retrieving the list of submitted URLs.
// It fetches the URLs from the storage service and returns them in the response.
func (h *Handler) UrlsHandler(c *gin.Context) {
	urls := services.GetSubmittedUrls()
	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
