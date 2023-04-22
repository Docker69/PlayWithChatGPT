package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	models "backend/models"
)

// get human by nickname
func (c *HumansCollectionType) GetByNickname(ctx context.Context, nickname string) (models.Human, error) {
	filter := bson.D{{Key: "nickname", Value: nickname}}
	result := c.col.FindOne(ctx, filter)

	var human models.Human
	if err := result.Decode(&human); err != nil {
		return human, err
	}

	return human, nil
}

// insert human in DB
func (c *HumansCollectionType) Insert(ctx context.Context, human *models.Human) (string, error) {
	insertResult, err := c.col.InsertOne(ctx, human)
	if err != nil {
		return "", err
	}

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// update human in DB by ID
func (c *HumansCollectionType) UpdateChats(ctx context.Context, human *models.Human) error {
	id, err := primitive.ObjectIDFromHex(human.Id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "chatids", Value: human.ChatIds}}}}

	_, err = c.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
