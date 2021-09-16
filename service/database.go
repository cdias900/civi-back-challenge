package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database service
type DBService interface {
	Connect(ctx context.Context, atlasURI string) error
	CreateDocument(ctx context.Context, collection string, input interface{}) error
	ReadDocuments(ctx context.Context, collection string, filter bson.M, limit int64, skip int64) (cursor *mongo.Cursor, err error)
	ReadDocument(ctx context.Context, collection string, filter bson.M) (output interface{}, err error)
	ReplaceDocument(ctx context.Context, collection string, filter bson.M, input interface{}) error
}

// Database service struct
type dbService struct {
	client *mongo.Client
}

// New database service
func NewDBService() DBService {
	return &dbService{
		client: nil,
	}
}

// Encode interface to BSON
func encodeBSON(v interface{}) *bson.M {
	bytes, err := bson.Marshal(v)
	if err != nil {
		log.Println("couldn't marshal interface into valid BSON byte array:", err)
		return nil
	}

	var b bson.M
	err = bson.Unmarshal(bytes, &b)
	if err != nil {
		log.Println("couldn't unmarshal byte array into valid BSON:", err)
		return nil
	}

	return &b
}

// Connect to database
func (dbS *dbService) Connect(ctx context.Context, atlasURI string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(atlasURI))
	if err != nil {
		return err
	}

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	dbS.client = client
	return nil
}

// Create document
func (dbS *dbService) CreateDocument(ctx context.Context, collection string, input interface{}) error {
	mainDatabase := dbS.client.Database("main")
	specificCollection := mainDatabase.Collection(collection)

	_, err := specificCollection.InsertOne(ctx, encodeBSON(input))
	return err
}

// Read many documents
func (dbS *dbService) ReadDocuments(ctx context.Context, collection string, filter bson.M, limit int64, skip int64) (cursor *mongo.Cursor, err error) {
	mainDatabase := dbS.client.Database("main")
	specificCollection := mainDatabase.Collection(collection)

	cur, err := specificCollection.Find(ctx, filter, &options.FindOptions{Limit: &limit, Skip: &skip})
	return cur, err
}

// Read document
func (dbS *dbService) ReadDocument(ctx context.Context, collection string, filter bson.M) (output interface{}, err error) {
	mainDatabase := dbS.client.Database("main")
	specificCollection := mainDatabase.Collection(collection)

	document := specificCollection.FindOne(ctx, filter)
	if document.Err() != nil {
		return nil, document.Err()
	}

	err = document.Decode(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Replace document
func (dbS *dbService) ReplaceDocument(ctx context.Context, collection string, filter bson.M, input interface{}) error {
	mainDatabase := dbS.client.Database("main")
	specificCollection := mainDatabase.Collection(collection)

	res := specificCollection.FindOneAndReplace(ctx, filter, encodeBSON(input))
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
