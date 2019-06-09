package storage_test

import (
	"context"
	"github.com/fernetbalboa/arqweb/src/config"
	"github.com/fernetbalboa/arqweb/src/mock"
	"github.com/fernetbalboa/arqweb/src/storage"
	"github.com/fernetbalboa/arqweb/src/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const (
	UserTestCollection = "user_test"
)

func init() {
	config.Scope = config.TestScope
}

func TestSaveUser(t *testing.T) {
	//Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userCollectionMock := mock.NewMockICollection(ctrl)

	documentId := test.NewDocumentId()

	userCollectionMock.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{
		InsertedID: documentId,
	}, nil)

	UserStorage, _ := storage.CreateUserStorage(userCollectionMock)

	user := test.DefaultUser()
	documentId = test.NewDocumentId()

	userCollectionMock.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{
		InsertedID: documentId,
	}, nil)

	//When
	savedUser, err := UserStorage.SaveUser(user)

	//Then
	assert.NoError(t, err)
	assert.Equal(t, documentId, savedUser.Id)
	assert.Equal(t, user.Username, savedUser.Username)
	assert.Equal(t, user.Password, savedUser.Password)
	assert.Equal(t, user.IsAdmin, savedUser.IsAdmin)
}

//Integration test using local MongoDB
func TestUserSearch(t *testing.T) {
	//Given
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	UserCollection := client.Database("test_db").Collection(UserTestCollection)

	UserStorage, _ := storage.CreateUserStorage(UserCollection)

	t.Run("Search by username", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database("test_db").Drop(context.Background())
		User := test.DefaultUser()
		anotherUser := test.DefaultUser()
		anotherUser.Username = "ALBERTO FERNANDEZ"

		savedUser, saveErr1 := UserStorage.SaveUser(User)
		anotherUsernameUser, saveErr2 := UserStorage.SaveUser(anotherUser)

		//When
		foundUser, searchErr := UserStorage.Search(User.Username)

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, searchErr)

		assert.Equal(t, foundUser.Username, savedUser.Username)
		assert.Equal(t, foundUser.Password, savedUser.Password)
		assert.Equal(t, foundUser.IsAdmin, savedUser.IsAdmin)
		assert.NotEqual(t, foundUser.Username, anotherUsernameUser.Username)
	})
}
