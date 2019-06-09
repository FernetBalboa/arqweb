package controller_test

import (
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/fernetbalboa/arqweb/config"
	"github.com/fernetbalboa/arqweb/controller"
	"github.com/fernetbalboa/arqweb/mock"
	"github.com/fernetbalboa/arqweb/test"
	"github.com/gin-gonic/gin/json"
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

func TestAddCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addCategoryEndpoint := "/categories"

	t.Run("Add category successfully", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockCategoryStorage(ctrl)
		CatController := &controller.CategoryController{CategoryStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addCategoryEndpoint, CatController.AddCategory)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addCategoryEndpoint, strings.NewReader(test.DefaultCategoryString()))

		savedCategory := test.DefaultCategoryWithId()
		storageMock.EXPECT().SaveCategory(test.DefaultCategory()).Return(savedCategory, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusCreated, w.Code)
		respBody, _ := json.Marshal(savedCategory)

		assert.Equal(t, respBody, w.Body.Bytes())

	})

	t.Run("Add category fails", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockCategoryStorage(ctrl)
		CatController := &controller.CategoryController{CategoryStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addCategoryEndpoint, CatController.AddCategory)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addCategoryEndpoint, strings.NewReader(test.DefaultCategoryString()))

		anError := apierror.New("anError")
		storageMock.EXPECT().SaveCategory(test.DefaultCategory()).Return(nil, anError)

		//When
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	t.Run("Add category bad request", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockCategoryStorage(ctrl)
		CatController := &controller.CategoryController{CategoryStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addCategoryEndpoint, CatController.AddCategory)
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", addCategoryEndpoint, strings.NewReader("badRequest"))

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
