package handlers

import (
	"github.com/gin-gonic/gin"
)

// Serve the embedded files via simple file serving; for now, read from disk.
func DocsHTML(c *gin.Context) {
	c.File("internal/docs/redoc.html")
}

func OpenAPIYAML(c *gin.Context) {
	c.File("internal/docs/openapi.yaml")
}

