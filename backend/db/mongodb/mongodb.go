package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend/utils"
)

var ChatsCollection *ChatsCollectionType = nil
var HumansCollection *HumansCollectionType = nil
var ConfigsCollection *ConfigsCollectionType = nil

var envVars = [...]string{"MONGO_HOST", "MONGO_PORT", "MONGO_USER_NAME", "MONGO_USER_PASSWORD", "MONGO_DATABASE"}

// create connection with mongo db
func init() {
	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Infof("Error loading .env file, will use default values. Err: %s", err)
	}

	// get env variables
	env := make(map[string]string)
	var exists bool = false
	for _, v := range envVars {
		env[v], exists = os.LookupEnv(v)
		if !exists {
			utils.Logger.Panic(v + " not found, panicking!!!")
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
		utils.Logger.Fatalf("Error connecting to MongoDB. Err: %s, connStr: %s", err, connStr)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		utils.Logger.Fatalf("Error pinging MongoDB. Err: %s", err)
	}

	utils.Logger.Info("Connected to MongoDB!")

	ChatsCollection = NewChatsCollection(client.Database(env["MONGO_DATABASE"]).Collection("chats"))
	HumansCollection = NewHumansCollection(client.Database(env["MONGO_DATABASE"]).Collection("humans"))
	ConfigsCollection = NewConfigsCollection(client.Database(env["MONGO_DATABASE"]).Collection("configs"))

	utils.Logger.Info("MongoDB collections initialized!")
}

// Disconnect from MongoDB
func Shutdown(ctx context.Context) error {

	//Close the connection to MongoDB from either collection
	if err := ChatsCollection.col.Database().Client().Disconnect(ctx); err != nil {
		utils.Logger.Errorf("Error disconnecting from MongoDB. Err: %s", err)
		return err
	}

	utils.Logger.Info("Connection to MongoDB closed.")
	return nil
}
