package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	completionmodels "backend/models"
	mylogger "backend/utils"
)

// collection object/instance
var chatsCollection *mongo.Collection

var envVars = [...]string{"MONGO_HOST", "MONGO_PORT", "MONGO_USER_NAME", "MONGO_USER_PASSWORD", "MONGO_DATABASE"}

// create connection with mongo db
func init() {
	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		mylogger.Logger.Panicf("Error loading .env file. Err: %s", err)
	}

	// get env variables
	env := make(map[string]string)
	var exists bool = false
	for _, v := range envVars {
		env[v], exists = os.LookupEnv(v)
		if !exists {
			mylogger.Logger.Panic(v + " not found, panicking!!!")
		}
	}
	//Build the connection string
	connStr := "mongodb://" + env["MONGO_USER_NAME"] + ":" + env["MONGO_USER_PASSWORD"] + "@" + env["MONGO_HOST"] + ":" + env["MONGO_PORT"] + "/" + env["MONGO_DATABASE"] + "?authSource=" + env["MONGO_DATABASE"]

	// Set client options
	clientOptions := options.Client().ApplyURI(connStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		mylogger.Logger.Fatalf("Error connecting to MongoDB. Err: %s, connStr: %s", err, connStr)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		mylogger.Logger.Fatalf("Error pinging MongoDB. Err: %s", err)
	}

	mylogger.Logger.Info("Connected to MongoDB!")

	chatsCollection = client.Database(env["MONGO_DATABASE"]).Collection("chats")

	mylogger.Logger.Info("chats collection instance created!")
}

// get all the chats from the DB and return them and error
func GetAllChats() ([]completionmodels.ChatCompletionRequestBody, error) {
	cur, err := chatsCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		mylogger.Logger.Fatalf("Error getting all chats. Err: %s", err)
	}
	defer cur.Close(context.Background())

	var results []completionmodels.ChatCompletionRequestBody
	for cur.Next(context.Background()) {
		var result completionmodels.ChatCompletionRequestBody
		e := cur.Decode(&result)
		if e != nil {
			mylogger.Logger.Fatalf("Error decoding chat. Err: %s", e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		mylogger.Logger.Fatalf("Error getting all chats. Err: %s", err)
	}

	mylogger.Logger.Debug("Get all chats successful!")

	//return results and error
	return results, err
}

// Insert new Chat in the DB
func InitNewChatDocument(chat *completionmodels.ChatCompletionRequestBody) (string, error) {
	insertResult, err := chatsCollection.InsertOne(context.Background(), chat)

	if err != nil {
		mylogger.Logger.Errorf("Error inserting chat. Err: %s", err)
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Update the chat in the DB
func UpdateChat(chat *completionmodels.ChatCompletionRequestBody) error {
	id, _ := primitive.ObjectIDFromHex(chat.Id)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"messages": chat.Messages}}

	_, err := chatsCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		mylogger.Logger.Errorf("Error updating chat. Err: %s", err)
	}

	return nil
}

// delete all the chats from the DB
func DeleteAllChats() int64 {
	d, err := chatsCollection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		mylogger.Logger.Fatalf("Error deleting all chats. Err: %s", err)
	}

	mylogger.Logger.Debugf("Deleted Document, count: %d", d.DeletedCount)
	return d.DeletedCount
}
