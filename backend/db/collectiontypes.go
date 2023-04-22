package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type ChatsCollectionType struct {
	col *mongo.Collection
}

func NewChatsCollection(col *mongo.Collection) *ChatsCollectionType {
	return &ChatsCollectionType{col: col}
}

type HumansCollectionType struct {
	col *mongo.Collection
}

func NewHumansCollection(col *mongo.Collection) *HumansCollectionType {
	return &HumansCollectionType{col: col}
}
