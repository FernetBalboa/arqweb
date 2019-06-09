package config

import (
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/gin-gonic/gin"
)

func ConfiguredRouter() *gin.Engine {
	router := gin.New()

	//Recover from panics and errors
	router.Use(apierror.CatchAPIErrors(), enableCORS())

	return router
}

//Enable CORS for Angular app
func enableCORS()gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
}