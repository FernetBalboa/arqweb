package storage_test

import (
	"context"
	"github.com/fernetbalboa/arqweb/src/domain"
	"github.com/fernetbalboa/arqweb/src/storage"
	"github.com/fernetbalboa/arqweb/src/test"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const (
	SuggestionTestCollection string = "sugg_test"
)

func TestGetSuggestions(t *testing.T) {
	//Given
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	SuggestionCollection := client.Database("test_db").Collection(SuggestionTestCollection)

	SuggestionStorage, _ := storage.CreateSuggestionStorage(SuggestionCollection)

	t.Run("Get pending suggestions", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database("test_db").Drop(context.Background())
		suggestion := test.DefaultCategorySuggestion()
		anotherSuggestion := test.DefaultCategorySuggestion()
		anotherSuggestion.CategoryName = "anotherName"
		anotherSuggestion.Status = domain.Approved
		anotherSuggestion2 := test.DefaultCategorySuggestion()
		anotherSuggestion2.CategoryName = "anotherName2"
		anotherSuggestion2.Status = domain.Rejected
		anotherSuggestion3 := test.DefaultCategorySuggestion()
		anotherSuggestion3.CategoryName = "anotherName3"
		anotherSuggestion3.Status = domain.Pending

		_, saveErr1 := SuggestionStorage.SaveSuggestion(suggestion)
		_, saveErr2 := SuggestionStorage.SaveSuggestion(anotherSuggestion)
		_, saveErr3 := SuggestionStorage.SaveSuggestion(anotherSuggestion2)
		_, saveErr4 := SuggestionStorage.SaveSuggestion(anotherSuggestion3)

		//When
		foundSuggestions, searchErr := SuggestionStorage.GetPendingSuggestions()

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, saveErr3)
		assert.NoError(t, saveErr4)
		assert.NoError(t, searchErr)

		assert.Contains(t, foundSuggestions, *suggestion)
		assert.NotContains(t, foundSuggestions, *anotherSuggestion)
		assert.NotContains(t, foundSuggestions, *anotherSuggestion2)
		assert.Contains(t, foundSuggestions, *anotherSuggestion3)
	})
}

func TestApproveOrRejectSuggestion(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	SuggestionCollection := client.Database("test_db").Collection(SuggestionTestCollection)

	SuggestionStorage, _ := storage.CreateSuggestionStorage(SuggestionCollection)

	t.Run("approve and reject suggestions", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database("test_db").Drop(context.Background())
		suggestion := test.DefaultCategorySuggestion()
		anotherSuggestion := test.DefaultCategorySuggestion()
		anotherSuggestion.CategoryName = "anotherName"
		anotherSuggestion2 := test.DefaultCategorySuggestion()
		anotherSuggestion2.CategoryName = "anotherName2"

		savedSuggestion, saveErr1 := SuggestionStorage.SaveSuggestion(suggestion)
		anotherSavedSuggestion, saveErr2 := SuggestionStorage.SaveSuggestion(anotherSuggestion)
		anotherSavedSuggestion2, saveErr3 := SuggestionStorage.SaveSuggestion(anotherSuggestion2)

		//When
		updatedSuggestion, approveErr := SuggestionStorage.ApproveOrRejectSuggestion(savedSuggestion.Id, true)
		updatedSuggestion2, approveErr2 := SuggestionStorage.ApproveOrRejectSuggestion(anotherSavedSuggestion2.Id, false)
		foundSuggestions, searchErr := SuggestionStorage.GetPendingSuggestions()

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, saveErr3)
		assert.NoError(t, searchErr)
		assert.NoError(t, approveErr)
		assert.NoError(t, approveErr2)

		assert.Equal(t, updatedSuggestion.Status, domain.Approved)
		assert.Equal(t, updatedSuggestion2.Status, domain.Rejected)
		assert.NotContains(t, foundSuggestions, savedSuggestion)
		assert.NotContains(t, foundSuggestions, anotherSavedSuggestion2)
		assert.Contains(t, foundSuggestions, *anotherSavedSuggestion)
	})
}
