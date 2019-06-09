package storage

import (
	"context"
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/fernetbalboa/arqweb/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

//go:generate mockgen -destination=../mock/mock_user_storage.go -package=mock -source=user_storage.go
type UserStorage interface {
	SaveUser(User *domain.User) (*domain.User, error)
	Search(userName string) (*domain.User, error)
}

const (
	UserCollection = "user"
)

func init() {
	//resetUserCollection() //Comment if data should be kept between program runs
}

type UserStorageImpl struct {
	UserCollection ICollection
}

func CreateUserStorage(UserCollection ICollection) (UserStorage, error) {
	storage := &UserStorageImpl{
		UserCollection: UserCollection,
	}
	_, err := storage.SaveUser(&domain.User{
		Username: "admin",
		Password: "admin",
		IsAdmin:  true,
	})
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func NewUserStorage() (UserStorage, error) {
	client, err := getMongoDBClient()
	UserCollection := client.Database(Database).Collection(UserCollection)
	if err != nil {
		return nil, err
	}

	return CreateUserStorage(UserCollection)
}

func (us *UserStorageImpl) SaveUser(User *domain.User) (*domain.User, error) {
	User.Id = primitive.NewObjectID()
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := us.UserCollection.InsertOne(ctx, User)

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not insert new User into MongoDB. User: %+v", User)
	}

	User.Id = res.InsertedID.(primitive.ObjectID)

	return User, nil
}

func (us *UserStorageImpl) Search(username string) (*domain.User, error) {

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	res := us.UserCollection.FindOne(ctx, bsonx.Doc{{"username", bsonx.String(username)}})
	var user domain.User
	err := res.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Errorf("Error while decoding User. USer: %v, doesnt exist.", username)
			return nil, err
		}
		log.Errorf("Error while decoding User. Cause: %v", err)
		return nil, err
	}
	return &user, nil
}

func resetUserCollection() {
	client, err := getMongoDBClient()

	if err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = client.Database(Database).Collection(UserCollection).Drop(ctx)

	if err != nil {
		log.Error("Could not reset User MongoDB collection")
	}
}
