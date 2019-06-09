package storage

import (
	"context"
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/fernetbalboa/arqweb/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

//go:generate mockgen -destination=../mock/mock_cat_suggestion_storage.go -package=mock -source=suggestion_storage.go
type SuggestionStorage interface {
	SaveSuggestion(suggestion *domain.CategorySuggestion) (*domain.CategorySuggestion, error)
	ApproveOrRejectSuggestion(id primitive.ObjectID, approve bool) (*domain.CategorySuggestion, error)
	GetPendingSuggestions() ([]domain.CategorySuggestion, error)
}

const (
	SuggestionCollection = "categorySuggestion"
)

func init() {
	resetSuggestionCollection() //Comment if data should be kept between program runs
}

type SuggestionStorageImpl struct {
	suggestionCollection ICollection
}

func CreateSuggestionStorage(suggestionCollection ICollection) (SuggestionStorage, error) {
	storage := &SuggestionStorageImpl{
		suggestionCollection: suggestionCollection,
	}
	return storage, nil
}

func NewSuggestionStorage() (SuggestionStorage, error) {
	client, err := getMongoDBClient()
	suggestionCollection := client.Database(Database).Collection(SuggestionCollection)
	if err != nil {
		return nil, err
	}

	return CreateSuggestionStorage(suggestionCollection)
}

func (ss *SuggestionStorageImpl) SaveSuggestion(suggestion *domain.CategorySuggestion) (*domain.CategorySuggestion, error) {
	suggestion.Id = primitive.NewObjectID()
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := ss.suggestionCollection.InsertOne(ctx, suggestion)

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not insert new suggestion into MongoDB. suggestion: %+v", suggestion)
	}

	suggestion.Id = res.InsertedID.(primitive.ObjectID)

	return suggestion, nil
}

func (ss *SuggestionStorageImpl) ApproveOrRejectSuggestion(id primitive.ObjectID, approve bool) (*domain.CategorySuggestion, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	log.Infof("approving or rejecting suggestion '%s'", id)

	status := domain.Rejected
	if approve {
		status = domain.Approved
	}
	res := ss.suggestionCollection.FindOne(ctx, bson.M{"_id": id})
	var suggestion domain.CategorySuggestion
	err := res.Decode(&suggestion)
	if err != nil {
		log.Errorf("Error while decoding suggestion. Cause: %v", err)
		return nil, err
	}
	suggestion.Status = status

	// set filters and updates
	filter := bson.M{"_id": id}
	update := bson.M{"$set": suggestion}

	_, err = ss.suggestionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, apierror.Wrapf(err, "Couldn't approve or reject suggestion with id: '%s'", id)
	}

	log.Infof("Suggestion '%s' successfully approved or rejected", id)
	return &suggestion, nil
}

func (ss *SuggestionStorageImpl) GetPendingSuggestions() ([]domain.CategorySuggestion, error) {

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	pendingString := string(domain.Pending)
	resCursor, err := ss.suggestionCollection.Find(ctx, bsonx.Doc{{"status", bsonx.String(pendingString)}})

	if err != nil {
		return []domain.CategorySuggestion{}, apierror.Wrapf(err, "Could not retrieve list of suggestions")
	}

	defer resCursor.Close(ctx)

	results := []domain.CategorySuggestion{}
	for resCursor.Next(ctx) {
		var sugg domain.CategorySuggestion
		err = resCursor.Decode(&sugg)
		if err != nil {
			log.Errorf("Error while decoding suggestion. Cause: %v", err)
		} else {
			results = append(results, sugg)
		}
	}

	return results, nil
}

func resetSuggestionCollection() {
	client, err := getMongoDBClient()

	if err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = client.Database(Database).Collection(SuggestionCollection).Drop(ctx)

	if err != nil {
		log.Error("Could not reset Suggestion MongoDB collection")
	}
}
