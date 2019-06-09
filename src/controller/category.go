package controller

import (
	"github.com/fernetbalboa/arqweb/src/apierror"
	"github.com/fernetbalboa/arqweb/src/domain"
	"github.com/fernetbalboa/arqweb/src/storage"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CategoryController struct {
	CategoryStorage storage.CategoryStorage
}

func CreateCategoryController(CategoryStorage storage.CategoryStorage) *CategoryController {
	return &CategoryController{
		CategoryStorage: CategoryStorage,
	}
}

func NewCategoryController() (*CategoryController, error) {
	CatStorage, err := storage.NewCategoryStorage()
	if err != nil {
		return nil, apierror.Wrapf(err, "Could not create Category controller")
	}

	return CreateCategoryController(CatStorage), nil
}

func (cc *CategoryController) AddCategory(c *gin.Context) {
	var category domain.Category
	err := c.ShouldBindJSON(&category)
	if err != nil || category.Name == "" {
		errorMsg := "Error parsing Category. It should be a Category obj"
		var apiError error
		if err != nil {
			apiError = apierror.BadRequest.Wrapf(err, errorMsg)
		} else {
			apiError = apierror.BadRequest.Newf(errorMsg)
		}
		_ = c.Error(apiError)
		return
	}

	savedCategory, err := cc.CategoryStorage.SaveCategory(&category)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, savedCategory)
}

func (cc *CategoryController) RemoveCategory(c *gin.Context) {
	categoryId := c.Param("id")

	res, err := cc.CategoryStorage.RemoveCategory(categoryId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (cc *CategoryController) GetCategories(c *gin.Context) {
	categories, err := cc.CategoryStorage.GetCategories()

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (cc *CategoryController) SearchCategory(c *gin.Context) {
	var searchFilters domain.CategoryFilter

	if err := c.ShouldBindQuery(&searchFilters); err != nil {
		err = apierror.BadRequest.Wrapf(err, "Invalid search query filters")
		_ = c.Error(err)
		return
	}

	if searchFilters.Limit == 0 {
		searchFilters.Limit = DefaultSearchLimit
	}

	log.Infof("Searching POIs for request: %s", c.Request.URL)

	categories, err := cc.CategoryStorage.SearchCategory(&searchFilters)

	if err != nil {
		_ = c.Error(err)
		return
	}

	log.Infof("Search results: %v", categories)
	c.JSON(http.StatusOK, categories)
}

func (cc *CategoryController) EditCategory(c *gin.Context) {
	categoryId := c.Param("id")
	var newVersionCategory domain.Category
	err := c.ShouldBindJSON(&newVersionCategory)

	if err != nil {
		_ = c.Error(err)
		return
	}

	if categoryId != newVersionCategory.Id.String() {
		apiError := apierror.BadRequest.New("Category field ID can not be updated")
		_ = c.Error(apiError)
		return
	}

	err = cc.CategoryStorage.EditCategory(&newVersionCategory)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": categoryId, "message": "Category successfully updated",
		"status": "OK", "code": http.StatusOK})
}
