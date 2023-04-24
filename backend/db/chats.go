package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	models "backend/models"
)

// get all chats from DB
func (c *ChatsCollectionType) GetAll(ctx context.Context) ([]models.ChatCompletionRequestBody, error) {
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
func (c *ChatsCollectionType) Insert(ctx context.Context, chat *models.ChatCompletionRequestBody) (string, error) {
	insertResult, err := c.col.InsertOne(ctx, chat)
	if err != nil {
		return "", err
	}

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

// get chat by ID from DB
func (c *ChatsCollectionType) GetById(ctx context.Context, id string) (models.ChatCompletionRequestBody, error) {
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
func (c *ChatsCollectionType) Update(ctx context.Context, chat *models.ChatCompletionRequestBody) error {
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
func (c *ChatsCollectionType) DeleteAll(ctx context.Context) (int64, error) {
	d, err := c.col.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		return 0, err
	}

	return d.DeletedCount, nil
}
