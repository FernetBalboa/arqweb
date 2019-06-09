package test

import (
	"encoding/json"
	"github.com/fernetbalboa/arqweb/src/domain"
)

func DefaultCategorySuggestion() *domain.CategorySuggestion {
	return &domain.CategorySuggestion{
		CategoryName:   "aName",
		HiddenCategory: false,
		CategoryIcon:   "anIcon",
		Status:         domain.Pending,
	}
}

func DefaultCategorySuggestionWithId() *domain.CategorySuggestion {
	return &domain.CategorySuggestion{
		Id:             NewDocumentId(),
		CategoryName:   "aName",
		HiddenCategory: false,
		CategoryIcon:   "anIcon",
		Status:         domain.Pending,
	}
}

func DefaultCategorySuggestionString() string {
	category := DefaultCategorySuggestion()
	out, err := json.Marshal(category)
	if err != nil {
		panic(err)
	}

	return string(out)
}
