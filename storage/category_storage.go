package storage

import (
	"context"
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/fernetbalboa/arqweb/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//go:generate mockgen -destination=../mock/mock_cat_storage.go -package=mock -source=category_storage.go
type CategoryStorage interface {
	SaveCategory(category *domain.Category) (*domain.Category, error)
	GetCategories() ([]domain.Category, error)
	SearchCategory(filters *domain.CategoryFilter) ([]*domain.Category, error)
	EditCategory(newVersionCategory *domain.Category) error
	RemoveCategory(categoryId string) (*mongo.DeleteResult, error)
}

const (
	CategoryCollection = "categories"
)

func init() {
	resetCategoryCollection() //Comment if data should be kept between program runs
}

type CategoryStorageImpl struct {
	catCollection ICollection
}

func CreateCategoryStorage(catCollection ICollection) (CategoryStorage, error) {
	storage := &CategoryStorageImpl{
		catCollection: catCollection,
	}

	_, err := storage.SaveCategory(NewCategory("food", "fas fa-hamburger"))
	if err != nil {
		return nil, err
	}

	_, err = storage.SaveCategory(NewCategory("art", "fas fa-utensils"))
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func NewCategoryStorage() (CategoryStorage, error) {
	client, err := getMongoDBClient()
	catCollection := client.Database(Database).Collection(CategoryCollection)
	if err != nil {
		return nil, err
	}

	return CreateCategoryStorage(catCollection)
}

func (cs *CategoryStorageImpl) SaveCategory(category *domain.Category) (*domain.Category, error) {

	category.Id = primitive.NewObjectID()
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := cs.catCollection.InsertOne(ctx, category)

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not insert new category into MongoDB. category: %+v", category)
	}

	category.Id = res.InsertedID.(primitive.ObjectID)

	return category, nil
}

func (cs *CategoryStorageImpl) RemoveCategory(categoryId string) (*mongo.DeleteResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	id, _ := primitive.ObjectIDFromHex(categoryId)
	res, err := cs.catCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not remove category from MongoDB. category: %s", categoryId)
	}

	return res, nil
}

func (cs *CategoryStorageImpl) GetCategories() ([]domain.Category, error) {

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	resCursor, err := cs.catCollection.Find(ctx, bson.M{})

	if err != nil {
		return []domain.Category{}, apierror.Wrapf(err, "Could not retrieve list of categories")
	}

	defer resCursor.Close(ctx)

	results := []domain.Category{}
	for resCursor.Next(ctx) {
		var cat domain.Category
		err = resCursor.Decode(&cat)
		if err != nil {
			log.Errorf("Error while decoding Category. Cause: %v", err)
		} else {
			results = append(results, cat)
		}
	}

	return results, nil
}

func (cs *CategoryStorageImpl) SearchCategory(filters *domain.CategoryFilter) ([]*domain.Category, error) {
	if filters == nil {
		filters = &domain.CategoryFilter{}
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	findOptions := options.Find()
	findOptions.SetLimit(filters.Limit)

	queryFilters := buildCategoryQueryFilters(filters)
	resCursor, err := cs.catCollection.Find(ctx, queryFilters, findOptions)

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not find any category using filters %+v", filters)
	}

	defer resCursor.Close(ctx)

	//noinspection GoPreferNilSlice
	results := []*domain.Category{}
	for resCursor.Next(ctx) {
		var cat domain.Category
		err = resCursor.Decode(&cat)
		if err != nil {
			log.Errorf("Error while decoding category. Cause: %v", err)
		} else {
			results = append(results, &cat)
		}
	}

	return results, nil

}

func (cs *CategoryStorageImpl) EditCategory(newVersionCategory *domain.Category) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	log.Infof("Updating category '%s'", newVersionCategory.Id)

	// set filters and updates
	filter := bson.M{"_id": newVersionCategory.Id}
	update := bson.M{"$set": newVersionCategory}

	// update document
	_, err := cs.catCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return apierror.Wrapf(err, "Couldn't update category with id: '%s'", newVersionCategory.Id)
	}

	log.Infof("Category '%s' successfully updated", newVersionCategory.Id)
	return nil
}

func buildCategoryQueryFilters(filters *domain.CategoryFilter) bson.M {
	filtersMap := bson.M{}
	if filters.Name != "" {
		filtersMap["name"] = filters.Name
	}
	filtersMap["hidden"] = filters.Hidden

	return filtersMap
}

func NewCategory(name string, icon string) *domain.Category {
	return &domain.Category{
		Name:   name,
		Hidden: false,
		Icon:   icon,
	}
}

func resetCategoryCollection() {
	client, err := getMongoDBClient()

	if err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = client.Database(Database).Collection(CategoryCollection).Drop(ctx)

	if err != nil {
		log.Error("Could not reset Category MongoDB collection")
	}
}
