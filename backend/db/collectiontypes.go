package mongodb

import "go.mongodb.org/mongo-driver/mongo"

// chats collection
type ChatsCollectionType struct {
	col *mongo.Collection
}

// create new chats collection
func NewChatsCollection(col *mongo.Collection) *ChatsCollectionType {
	return &ChatsCollectionType{col: col}
}

// humans collection
type HumansCollectionType struct {
	col *mongo.Collection
}

// create new humans collection
func NewHumansCollection(col *mongo.Collection) *HumansCollectionType {
	return &HumansCollectionType{col: col}
}

// configs collection
type ConfigsCollectionType struct {
	col *mongo.Collection
}

// create new configs collection
func NewConfigsCollection(col *mongo.Collection) *ConfigsCollectionType {
	return &ConfigsCollectionType{col: col}
}
