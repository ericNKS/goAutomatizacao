package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth(c *gin.Context) {
	if false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
	}
	c.Next()
}
