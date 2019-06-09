package controller_test

import (
	"encoding/json"
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

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addUserEndpoint := "/user/signup"

	t.Run("Add User successfully", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockUserStorage(ctrl)
		UserController := &controller.UserController{UserStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addUserEndpoint, UserController.Signup)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addUserEndpoint, strings.NewReader(test.DefaultUserString()))

		savedUser := test.DefaultUserWithId()
		storageMock.EXPECT().Search(savedUser.Username).Return(nil, nil)
		storageMock.EXPECT().SaveUser(test.DefaultUser()).Return(savedUser, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusCreated, w.Code)
		respBody, _ := json.Marshal(savedUser)

		assert.Equal(t, respBody, w.Body.Bytes())

	})

	t.Run("Add already created user", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockUserStorage(ctrl)
		UserController := &controller.UserController{UserStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addUserEndpoint, UserController.Signup)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addUserEndpoint, strings.NewReader(test.DefaultUserString()))

		savedUser := test.DefaultUserWithId()
		storageMock.EXPECT().Search(savedUser.Username).Return(savedUser, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, `{"message":"There was a problem with the request","error":"forbidden","status":403,"cause":["Username already exists"]}`,
			w.Body.String())

	})
}

func TestSearchUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	loginEndpoint := "/user/login"

	t.Run("Login successfully", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockUserStorage(ctrl)
		UserController := &controller.UserController{UserStorage: storageMock}

		r := config.ConfiguredRouter()
		r.GET(loginEndpoint, UserController.Login)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", loginEndpoint, strings.NewReader(test.DefaultUserString()))

		savedUser := test.DefaultUser()

		storageMock.EXPECT().Search(savedUser.Username).Return(savedUser, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusOK, w.Code)
		respBody, _ := json.Marshal(savedUser)

		assert.Equal(t, respBody, w.Body.Bytes())

	})

	t.Run("login fails due to use an invalid password", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockUserStorage(ctrl)
		UserController := &controller.UserController{UserStorage: storageMock}

		r := config.ConfiguredRouter()
		r.GET(loginEndpoint, UserController.Login)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", loginEndpoint, strings.NewReader(test.DefaultUserString()))

		savedUser := test.DefaultUser()
		otherUser := savedUser
		otherUser.Password = "RICARDO FORT"
		storageMock.EXPECT().Search(savedUser.Username).Return(otherUser, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, `{"message":"There was a problem with the request","error":"forbidden","status":403,"cause":["Invalid username or password"]}`,
			w.Body.String())

	})

	t.Run("login fails: invalid user", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockUserStorage(ctrl)
		UserController := &controller.UserController{UserStorage: storageMock}

		r := config.ConfiguredRouter()
		r.GET(loginEndpoint, UserController.Login)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", loginEndpoint, strings.NewReader(test.DefaultUserString()))

		savedUser := test.DefaultUser()
		storageMock.EXPECT().Search(savedUser.Username).Return(nil, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, `{"message":"There was a problem with the request","error":"forbidden","status":403,"cause":["Invalid username or password"]}`,
			w.Body.String())

	})
}
