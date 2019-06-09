package config

import (
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ConfiguredRouter() *gin.Engine {
	router := gin.New()

	//Enable CORS for Angular app
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://arqweb-frontend.herokuapp.com"}

	//Recover from panics and errors
	router.Use(apierror.CatchAPIErrors(), cors.New(config))

	return router
}