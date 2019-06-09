package test

import (
	"github.com/fernetbalboa/arqweb/domain"
	"github.com/gin-gonic/gin/json"
)

func DefaultCategory() *domain.Category {
	return &domain.Category{
		Name:   "aName",
		Hidden: false,
		Icon:   "anIcon",
	}
}

func DefaultCategoryWithId() *domain.Category {
	return &domain.Category{
		Id:     NewDocumentId(),
		Name:   "aName",
		Hidden: false,
		Icon:   "anIcon",
	}
}

func DefaultCategoryString() string {
	category := DefaultCategory()
	out, err := json.Marshal(category)
	if err != nil {
		panic(err)
	}

	return string(out)
}
