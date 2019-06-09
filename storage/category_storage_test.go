package storage_test

import (
	"context"
	"github.com/fernetbalboa/arqweb/domain"
	"github.com/fernetbalboa/arqweb/storage"
	"github.com/fernetbalboa/arqweb/test"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestSearchCategory(t *testing.T) {
	//Given
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	CatCollection := client.Database(TestDB).Collection(CategoryTestCollection)

	POIStorage, _ := storage.CreateCategoryStorage(CatCollection)

	t.Run("Search category by name", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database(TestDB).Drop(context.Background())
		category := test.DefaultCategory()
		anotherCategory := test.DefaultCategory()
		anotherCategory.Name = "anotherName"

		_, saveErr1 := POIStorage.SaveCategory(category)
		_, saveErr2 := POIStorage.SaveCategory(anotherCategory)

		filter := &domain.CategoryFilter{Name: category.Name}

		//When
		foundCategories, searchErr := POIStorage.SearchCategory(filter)

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, searchErr)

		assert.Contains(t, foundCategories, category)
		assert.NotContains(t, foundCategories, anotherCategory)
	})

	t.Run("Search category by hidden", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database(TestDB).Drop(context.Background())
		category := test.DefaultCategory()
		hiddenCategory := test.DefaultCategory()
		hiddenCategory.Hidden = true

		_, saveErr1 := POIStorage.SaveCategory(category)
		_, saveErr2 := POIStorage.SaveCategory(hiddenCategory)

		filter := &domain.CategoryFilter{Hidden: false}

		//When
		foundCategories, searchErr := POIStorage.SearchCategory(filter)

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, searchErr)

		assert.Contains(t, foundCategories, category)
		assert.NotContains(t, foundCategories, hiddenCategory)
	})

}
