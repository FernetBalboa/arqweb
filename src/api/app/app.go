// Package app is the start package of CX Emails
// It has router start configuration and URL route mapping.
// Gin-Gonic Framework: https://github.com/gin-gonic/gin
// is used for routing and server working.
package app

import (
	"github.com/fernetbalboa/arqweb/src/api/config"
	"github.com/gin-gonic/gin"
	"os"
)

// Router is the base Gin Router used for mapping endpoints.
var Router *gin.Engine

// Start app point.
func Start() {
	Router = config.ConfiguredRouter()
	LoadEndpoints()

	port := os.Getenv("PORT")
	_ = Router.Run(":" + port)
}
