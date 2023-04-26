package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	models "backend/models"
)

// get human by nickname
func (c *HumansCollectionType) GetByNickname(nickname string) (models.Human, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "nickname", Value: nickname}}
	result := c.col.FindOne(ctx, filter)

	var human models.Human
	if err := result.Decode(&human); err != nil {
		return human, err
	}

	return human, nil
}

// insert human in DB
func (c *HumansCollectionType) Insert(human *models.Human) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := c.col.InsertOne(ctx, human)
	if err != nil {
		return "", err
	}

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// update human in DB by ID
func (c *HumansCollectionType) UpdateChats(human *models.Human) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

// get human by ID from DB
func (c *HumansCollectionType) GetById(id string) (models.Human, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Human{}, err
	}

	filter := bson.M{"_id": oid}
	result := c.col.FindOne(ctx, filter)

	var human models.Human
	if err := result.Decode(&human); err != nil {
		return models.Human{}, err
	}

	return human, nil
}

// get human by chat id
func (c *HumansCollectionType) GetByChatId(chatId string) (models.Human, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "chatids.id", Value: chatId}}
	result := c.col.FindOne(ctx, filter)

	var human models.Human
	if err := result.Decode(&human); err != nil {
		return human, err
	}

	return human, nil
}
