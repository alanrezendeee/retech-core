package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/auth"
)

func Me(c *gin.Context) {
	userID, _ := c.Get(auth.CtxUserID)
	roles, _ := c.Get(auth.CtxRoles)
	c.JSON(http.StatusOK, gin.H{"userId": userID, "roles": roles})
}

