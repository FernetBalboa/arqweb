package controller

import (
	"github.com/fernetbalboa/arqweb/src/apierror"
	"github.com/fernetbalboa/arqweb/src/domain"
	"github.com/fernetbalboa/arqweb/src/storage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type SuggestionController struct {
	SuggestionStorage storage.SuggestionStorage
	CategoryStorage   storage.CategoryStorage
}

func CreateSuggestionController(SuggestionStorage storage.SuggestionStorage, CategoryStorage storage.CategoryStorage) *SuggestionController {
	return &SuggestionController{
		SuggestionStorage: SuggestionStorage,
		CategoryStorage:   CategoryStorage,
	}
}

func NewSuggestionController(categoryStorage storage.CategoryStorage) (*SuggestionController, error) {
	SuggestionStorage, err := storage.NewSuggestionStorage()
	if err != nil {
		return nil, apierror.Wrapf(err, "Could not create Suggestion controller")
	}

	return CreateSuggestionController(SuggestionStorage, categoryStorage), nil
}

func (sc *SuggestionController) AddSuggestion(c *gin.Context) {
	var suggestion domain.CategorySuggestion
	err := c.ShouldBindJSON(&suggestion)
	if err != nil || suggestion.CategoryName == "" {
		errorMsg := "Error parsing Suggestion. It should be a Suggestion obj"
		var apiError error
		if err != nil {
			apiError = apierror.BadRequest.Wrapf(err, errorMsg)
		} else {
			apiError = apierror.BadRequest.Newf(errorMsg)
		}
		_ = c.Error(apiError)
		return
	}
	suggestion.Status = domain.Pending
	savedSuggestion, err := sc.SuggestionStorage.SaveSuggestion(&suggestion)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, savedSuggestion)
}

func (sc *SuggestionController) ApproveSuggestion(c *gin.Context) {
	id := c.Param("id")
	suggestionId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorMsg := "Error parsing SuggestionId. It should be a valid objectId"
		apiError := apierror.BadRequest.Wrapf(err, errorMsg)
		_ = c.Error(apiError)
		return
	}
	suggestion, err := sc.SuggestionStorage.ApproveOrRejectSuggestion(suggestionId, true)
	if err != nil {
		_ = c.Error(err)
		return
	}

	category := suggestionToCategory(suggestion)
	category, err = sc.CategoryStorage.SaveCategory(category)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (sc *SuggestionController) RejectSuggestion(c *gin.Context) {
	id := c.Param("id")
	suggestionId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorMsg := "Error parsing SuggestionId. It should be a valid objectId"
		apiError := apierror.BadRequest.Wrapf(err, errorMsg)
		_ = c.Error(apiError)
		return
	}

	_, err = sc.SuggestionStorage.ApproveOrRejectSuggestion(suggestionId, false)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": suggestionId, "message": "Suggestion successfully rejected",
		"status": "OK", "code": http.StatusOK})
}

func (sc *SuggestionController) GetSuggestions(c *gin.Context) {
	suggestions, err := sc.SuggestionStorage.GetPendingSuggestions()

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

func suggestionToCategory(suggestion *domain.CategorySuggestion) *domain.Category {
	return &domain.Category{
		Name:   suggestion.CategoryName,
		Hidden: suggestion.HiddenCategory,
		Icon:   suggestion.CategoryIcon,
	}
}
