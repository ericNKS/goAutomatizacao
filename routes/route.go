package routes

import (
	"automatizacao/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Route(app *gin.Engine) {
	app.GET("/server/:port", func(c *gin.Context) {
		port := c.Param("port")
		if len(port) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid port"})
			c.Abort()
		}
		if handler.IsGabOn(port) {
			c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
			c.AbortWithStatus(200)
		} else {
			c.JSON(http.StatusConflict, gin.H{"message": "Server is offline"})
		}

	})

	// Protected route
	// Route to open prompt
	//app.GET("/server/a/:dir/:port", middleware.Auth, handler.RunGab)
}
