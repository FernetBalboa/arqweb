package storage

import (
	"context"
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/fernetbalboa/arqweb/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const Database = "arqweb"

//go:generate mockgen -destination=../mock/mock_db_collection.go -package=mock -source=db.go
type ICollection interface {
	BulkWrite(ctx context.Context, models []mongo.WriteModel,
		opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)

	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	Find(ctx context.Context, filter interface{},
		opts ...*options.FindOptions) (*mongo.Cursor, error)

	FindOne(ctx context.Context, filter interface{},
		opts ...*options.FindOneOptions) *mongo.SingleResult

	FindOneAndUpdate(ctx context.Context, filter interface{},
		update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult

	DeleteOne(ctx context.Context, filter interface{},
		opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

func getMongoDBClient() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var dbUri string;

	if config.InDevelopment() {
		dbUri = "mongodb://localhost:27017"
	} else {
		dbUri = "mongodb://arqweb:arqweb123@ds263146.mlab.com:63146/arqweb"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))

	if err != nil {
		return nil, apierror.Wrap(err, "Could not connect to MongoDB")
	}
	log.Info("Succesfully connected to MonGoDB database")

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, apierror.Wrapf(err, "MongoDB client did not respond to ping")
	}

	return client, nil
}
