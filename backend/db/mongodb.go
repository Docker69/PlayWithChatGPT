package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mylogger "backend/utils"
)

var ChatsCollection *ChatsCollectionType = nil
var HumansCollection *HumansCollectionType = nil

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

	ChatsCollection = NewChatsCollection(client.Database(env["MONGO_DATABASE"]).Collection("chats"))
	HumansCollection = NewHumansCollection(client.Database(env["MONGO_DATABASE"]).Collection("humans"))

	mylogger.Logger.Info("MongoDB collections initialized!")
}
