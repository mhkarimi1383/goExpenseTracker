package database

import (
	"context"
	"strings"

	"github.com/mhkarimi1383/goExpenseTracker/configuration"
	"github.com/mhkarimi1383/goExpenseTracker/logger"
	"github.com/mhkarimi1383/goExpenseTracker/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongoClient  *mongo.Client
	ctx          context.Context
	databaseName string
	database     *mongo.Database
)

func init() {
	cfg, err := configuration.GetConfig()
	if err != nil {
		logger.Fatalf(true, "error in initializing configuration: %v", err)
	}
	ctx = context.Background()
	logger.Infof(false, "initializing MongoDB Client...")
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBConnectionURI))
	if err != nil {
		logger.Fatalf(true, "error while connecting to MongoDB: %v", err)
	}
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		logger.Fatalf(true, "error while pinging Primary server in MongoDB: %v", err)
	}
	logger.Infof(false, "Successfully connected to MongoDB :)")
	databaseName = strings.ReplaceAll(cfg.ApplicationTitle, " ", "")
	logger.Infof(false, "using database %v to store data in mongoDB", databaseName)
	database = mongoClient.Database(databaseName)
}

func ListDatabaseNames() ([]string, error) {
	return mongoClient.ListDatabaseNames(ctx, bson.M{})
}

func InsertItem(username string, item types.Item) (*mongo.InsertOneResult, error) {
	collection := database.Collection(username)
	return collection.InsertOne(context.TODO(), item)
}

func ListItems(username string) ([]types.Item, error) {
	cur, currErr := database.Collection(username).Find(context.TODO(), bson.D{})
	if currErr != nil {
		return nil, currErr
	}
	defer cur.Close(ctx)

	var items []types.Item
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func DeleteItem(username string, id uint) (*mongo.DeleteResult, error) {
	return database.Collection(username).DeleteOne(context.TODO(), bson.M{"_id": id})
}
