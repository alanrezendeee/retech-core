package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	want := make(map[string]struct{}, len(roles))
	for _, r := range roles { want[r] = struct{}{} }
	return func(c *gin.Context) {
		got, ok := c.Get(CtxRoles)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"no roles"})
			return
		}
		okRole := false
		for _, r := range got.([]string) {
			if _, ok := want[r]; ok { okRole = true; break }
		}
		if !okRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"insufficient role"})
			return
		}
		c.Next()
	}
}

