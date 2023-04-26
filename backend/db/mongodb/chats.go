package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend/models"
)

// get all chats from DB
func (c *ChatsCollectionType) GetAll() ([]models.ChatCompletionRequestBody, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := c.col.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.ChatCompletionRequestBody
	for cur.Next(ctx) {
		var result models.ChatCompletionRequestBody
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// insert new chat in DB
func (c *ChatsCollectionType) Insert(chat *models.ChatCompletionRequestBody) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := c.col.InsertOne(ctx, chat)
	if err != nil {
		return "", err
	}

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// get chat by ID from DB
func (c *ChatsCollectionType) GetById(id string) (models.ChatCompletionRequestBody, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ChatCompletionRequestBody{}, err
	}

	filter := bson.M{"_id": oid}
	result := c.col.FindOne(ctx, filter)

	var chat models.ChatCompletionRequestBody
	if err := result.Decode(&chat); err != nil {
		return models.ChatCompletionRequestBody{}, err
	}

	return chat, nil
}

// update chat in DB
func (c *ChatsCollectionType) Update(chat *models.ChatCompletionRequestBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(chat.Id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"messages": chat.Messages}}

	_, err = c.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// delete all chats from DB
func (c *ChatsCollectionType) DeleteAll() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d, err := c.col.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		return 0, err
	}

	return d.DeletedCount, nil
}
