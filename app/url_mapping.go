package app

import (
	"github.com/fernetbalboa/arqweb/src/api/controller"
	log "github.com/sirupsen/logrus"
)

var poiController *controller.POIController
var categoryController *controller.CategoryController
var userController *controller.UserController
var suggestionController *controller.SuggestionController
var externalPOIController *controller.ExternalPOIController

func init() {
	CategoryController, err := controller.NewCategoryController()
	if err != nil {
		log.Fatalf("Could not create Category controller. Cause: %s", err.Error())
	}

	POIController, err := controller.NewPOIController(CategoryController.CategoryStorage)
	if err != nil {
		log.Fatalf("Could not create POI controller. Cause: %s", err.Error())
	}

	UserController, err := controller.NewUserController()
	if err != nil {
		log.Fatalf("Could not create User controller. Cause: %s", err.Error())
	}
	SuggestionController, err := controller.NewSuggestionController(CategoryController.CategoryStorage)
	if err != nil {
		log.Fatalf("Could not create Suggestion controller. Cause: %s", err.Error())
	}
	ExternalPOIController, err := controller.NewExternalPOIController(CategoryController.CategoryStorage)
	if err != nil {
		log.Fatalf("Could not create External POI controller. Cause: %s", err.Error())
	}
	poiController = POIController
	userController = UserController
	categoryController = CategoryController
	suggestionController = SuggestionController
	externalPOIController = ExternalPOIController
}

// LoadEndpoints is the base function to map endpoints.
func LoadEndpoints() {
	Router.GET("/ping", controller.Ping)

	userGroup := Router.Group("/user")
	userGroup.POST("/signup", userController.Signup)
	userGroup.POST("/login", userController.Login)

	poiGroup := Router.Group("/poi")

	poiGroup.POST("", poiController.AddPOI)
	poiGroup.GET("/search", poiController.SearchPOI)
	poiGroup.PUT("/:id", poiController.EditPOI)

	categoriesGroup := Router.Group("/categories")
	categoriesGroup.GET("", categoryController.GetCategories)
	categoriesGroup.GET("/search", categoryController.SearchCategory)
	categoriesGroup.POST("", categoryController.AddCategory)
	categoriesGroup.DELETE("/:id", categoryController.RemoveCategory)
	categoriesGroup.PUT("/:id", categoryController.EditCategory)

	suggestionGroup := Router.Group("/suggestions/categories/new")
	suggestionGroup.POST("", suggestionController.AddSuggestion)
	suggestionGroup.PUT("/:id/approve", suggestionController.ApproveSuggestion)
	suggestionGroup.PUT("/:id/reject", suggestionController.RejectSuggestion)
	suggestionGroup.GET("", suggestionController.GetSuggestions)


	externalGroup := Router.Group("/external")
	extPoiGroup := externalGroup.Group("/poi")

	extPoiGroup.POST("", externalPOIController.AddPOI)
	extPoiGroup.DELETE("/:id", externalPOIController.RemovePOI)
	extPoiGroup.GET("", externalPOIController.GetPOIs)
	extPoiGroup.GET("/:id", externalPOIController.GetPOI)
}
