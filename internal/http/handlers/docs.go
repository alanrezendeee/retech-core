package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Serve the embedded files via simple file serving; for now, read from disk.
func DocsHTML(c *gin.Context) {
	// Tentar caminho relativo primeiro, depois absoluto (para Docker)
	paths := []string{"internal/docs/redoc.html", "/app/internal/docs/redoc.html"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			http.ServeFile(c.Writer, c.Request, path)
			return
		}
	}
	c.String(404, "Documentação não encontrada")
}

func OpenAPIYAML(c *gin.Context) {
	// Tentar caminho relativo primeiro, depois absoluto (para Docker)
	var content []byte
	var err error

	paths := []string{"internal/docs/openapi.yaml", "/app/internal/docs/openapi.yaml"}
	for _, path := range paths {
		content, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		c.String(500, "Erro ao carregar OpenAPI spec: "+err.Error())
		return
	}

	// Substituir URL do servidor dinamicamente
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://api-core.theretech.com.br"
	}

	// Substituir o placeholder __API_BASE_URL__ no YAML
	yamlString := string(content)
	yamlString = strings.Replace(
		yamlString,
		"__API_BASE_URL__",
		apiBaseURL,
		-1, // Substituir todas as ocorrências
	)

	// Retornar YAML com URL dinâmica
	c.Header("Content-Type", "application/x-yaml")
	c.String(200, yamlString)
}
