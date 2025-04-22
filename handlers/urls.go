package handlers

import (
	"net/http"
	services "web-analyzer/internal/storage"

	"github.com/gin-gonic/gin"
)

func UrlsHandler(c *gin.Context) {
	urls := services.GetSubmittedUrls()

	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
