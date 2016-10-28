package api

import (
	"github.com/gin-gonic/gin"
)

func apiVersion(c *gin.Context) {
	c.JSON(200, gin.H{"Version": "1.0"})
}
