package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	AppName    = "retech-core"
	AppVersion = "0.1.0"
	GitCommit  = "dev"
	BuildDate  = "dev"
)

func Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":      AppName,
		"version":   AppVersion,
		"commit":    GitCommit,
		"buildDate": BuildDate,
	})
}

