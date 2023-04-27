package mongodb

import "go.mongodb.org/mongo-driver/mongo"

// chats collection
type ChatsCollectionType struct {
	col *mongo.Collection
}

// create new chats collection
func newChatsCollection(col *mongo.Collection) *ChatsCollectionType {
	return &ChatsCollectionType{col: col}
}

// humans collection
type HumansCollectionType struct {
	col *mongo.Collection
}

// create new humans collection
func newHumansCollection(col *mongo.Collection) *HumansCollectionType {
	return &HumansCollectionType{col: col}
}

// configs collection
type ConfigsCollectionType struct {
	col *mongo.Collection
}

// create new configs collection
func newConfigsCollection(col *mongo.Collection) *ConfigsCollectionType {
	return &ConfigsCollectionType{col: col}
}

// templates collection

type TemplatesCollectionType struct {
	col *mongo.Collection
}

// create new templates collection
func newTemplatesCollection(col *mongo.Collection) *TemplatesCollectionType {
	return &TemplatesCollectionType{col: col}
}

// AutoAIs model collection type
type AutoAIsCollectionType struct {
	col *mongo.Collection
}

// create new AutoAIs collection
func newAutoAIsCollection(col *mongo.Collection) *AutoAIsCollectionType {
	return &AutoAIsCollectionType{col: col}
}
