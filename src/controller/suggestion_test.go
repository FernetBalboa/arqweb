package controller_test

import (
	"encoding/json"
	"github.com/fernetbalboa/arqweb/src/apierror"
	"github.com/fernetbalboa/arqweb/src/config"
	"github.com/fernetbalboa/arqweb/src/controller"
	"github.com/fernetbalboa/arqweb/src/mock"
	"github.com/fernetbalboa/arqweb/src/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	config.Scope = config.TestScope
}

func TestAddCategorySuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addSuggestionEndpoint := "/suggestion/new"

	t.Run("Add suggestion successfully", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockSuggestionStorage(ctrl)
		categoryStorageMock := mock.NewMockCategoryStorage(ctrl)
		SuggestionController := &controller.SuggestionController{CategoryStorage: categoryStorageMock, SuggestionStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addSuggestionEndpoint, SuggestionController.AddSuggestion)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addSuggestionEndpoint, strings.NewReader(test.DefaultCategorySuggestionString()))

		savedSuggestion := test.DefaultCategorySuggestionWithId()
		storageMock.EXPECT().SaveSuggestion(test.DefaultCategorySuggestion()).Return(savedSuggestion, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusCreated, w.Code)
		respBody, _ := json.Marshal(savedSuggestion)

		assert.Equal(t, respBody, w.Body.Bytes())

	})

	t.Run("Add suggestion fails", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockSuggestionStorage(ctrl)
		categoryStorageMock := mock.NewMockCategoryStorage(ctrl)
		SuggestionController := &controller.SuggestionController{CategoryStorage: categoryStorageMock, SuggestionStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addSuggestionEndpoint, SuggestionController.AddSuggestion)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addSuggestionEndpoint, strings.NewReader(test.DefaultCategorySuggestionString()))

		anError := apierror.New("anError")
		storageMock.EXPECT().SaveSuggestion(test.DefaultCategorySuggestion()).Return(nil, anError)

		//When
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	t.Run("Add suggestion bad request", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockSuggestionStorage(ctrl)
		categoryStorageMock := mock.NewMockCategoryStorage(ctrl)
		SuggestionController := &controller.SuggestionController{CategoryStorage: categoryStorageMock, SuggestionStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addSuggestionEndpoint, SuggestionController.AddSuggestion)
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", addSuggestionEndpoint, strings.NewReader("badRequest"))

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
