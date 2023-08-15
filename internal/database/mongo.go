package database

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Client        *mongo.Client
	Database      *mongo.Database
	connectionStr string
	databaseName  string
}

func NewMongoConnection(connectionString, dbName string) *MongoConnection {
	return &MongoConnection{
		connectionStr: connectionString,
		databaseName:  dbName,
	}
}

func (mc *MongoConnection) Connect() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mc.connectionStr))
	if err != nil {
		return err
	}

	mc.Client = client
	mc.Database = client.Database(mc.databaseName)

	return nil
}

func (mc *MongoConnection) getCollection(collectionName string) *mongo.Collection {
	return mc.Database.Collection(collectionName)
}

func (mc *MongoConnection) Disconnect(ctx context.Context) {
	if err := mc.Client.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from database: %v", err)
	}
}

func (mc *MongoConnection) Insert(ctx context.Context, collectionName string, document interface{}) (string, error) {
	collection := mc.getCollection(collectionName)
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	return id.Hex(), nil
}

func (mc *MongoConnection) Exists(ctx context.Context, collectionName string, filters map[string]interface{}) (bool, error) {
	collection := mc.getCollection(collectionName)

	filter := bson.M{}
	for key, value := range filters {
		filter[key] = value
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (mc *MongoConnection) GetAll(ctx context.Context, collectionName string, filter map[string]interface{}) ([]interface{}, error) {
	collection := mc.getCollection(collectionName)

	// Преобразуем фильтр в bson.M, что на самом деле является псевдонимом для map[string]interface{}
	bsonFilter := bson.M(filter)

	cursor, err := collection.Find(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []interface{}
	for cursor.Next(ctx) {
		var item interface{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
