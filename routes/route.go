package routes

import (
	"automatizacao/handler"
	"automatizacao/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Route(app *gin.Engine) {
	app.GET("/server/:port", func(c *gin.Context) {
		port := c.Param("port")
		if len(port) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid port"})
		}
		if handler.IsGabOn(port) {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		}
		c.JSON(http.StatusConflict, gin.H{"error": "Port in use or is not an gabinete"})
	})

	// Protected route
	app.GET("/server/a/:dir/:port", middleware.Auth, handler.RunGab)
}
