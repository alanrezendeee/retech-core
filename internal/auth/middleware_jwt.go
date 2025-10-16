package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "userId"
const CtxRoles  = "roles"

func AuthJWT(jwt *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"missing bearer"})
			return
		}
		tok := strings.TrimPrefix(h, "Bearer ")
		_, claims, err := jwt.Parse(tok)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid token"})
			return
		}
		c.Set(CtxUserID, claims["sub"])
		if v, ok := claims["roles"].([]any); ok {
			rs := make([]string, 0, len(v))
			for _, e := range v { if s, ok := e.(string); ok { rs = append(rs, s) } }
			c.Set(CtxRoles, rs)
		}
		c.Next()
	}
}

