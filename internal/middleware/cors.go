package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allow bool, origins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !allow {
			c.Next()
			return
		}
		origin := c.GetHeader("Origin")
		allowed := "*"
		if len(origins) > 0 {
			for _, o := range origins {
				if strings.EqualFold(o, origin) {
					allowed = origin
					break
				}
			}
		}
		c.Header("Access-Control-Allow-Origin", allowed)
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-API-Key, X-Request-ID")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

