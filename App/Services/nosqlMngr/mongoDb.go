package nosqlmngr

import (
	"context"
	"fmt"
	"time"

	"github.com/Ahmed-AbdElRhman/twitter-architecture/tweets/tweets_services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	Client *mongo.Client
}

func NewMongoDb(uri string) (*MongoDb, error) {
	client, err := connectMongoDB(uri)
	if err != nil || client == nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoDb{
		Client: client,
	}, nil
}

// CreateTweet inserts a new tweet into the database.
func (obj *MongoDb) CreateTweet(tweet *tweets_services.Tweet, dbname, collectionName string) error {
	tweet.CreatedAt = time.Now()
	rslt, err := obj.getCollection(dbname, collectionName).InsertOne(context.Background(), tweet)
	fmt.Println("CreateTweet", rslt)
	return err
}

// GetUserTweets returns all tweets for a given user.
func (obj *MongoDb) GetUserTweets(ownerId *int, dbname, collectionName string) ([]*tweets_services.Tweet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if ownerId != nil {
		filter["owner"] = *ownerId
	}

	cursor, err := obj.getCollection(dbname, collectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tweets []*tweets_services.Tweet
	if err := cursor.All(ctx, &tweets); err != nil {
		return nil, err
	}
	fmt.Println("GetUserTweets", tweets)
	return tweets, nil
}

// CloseDB closes the connection to the database.
func (obj *MongoDb) CloseDB() error {
	return obj.Client.Disconnect(context.Background())
}

// Private and helper methods

// Get collection reference
func (obj *MongoDb) getCollection(dbname string, collectionName string) *mongo.Collection {
	return obj.Client.Database(dbname).Collection(collectionName)
}

// ConnectMongoDB establishes a connection to MongoDB and returns a collection reference.
func connectMongoDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
