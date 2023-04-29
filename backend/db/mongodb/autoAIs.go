package mongodb

import (
	"backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//implements get, insert, update, delete in the mongoDB for AutoAI model

// get all AutoAI models from AutoAIsCollection
func (c *AutoAIsCollectionType) GetAll() ([]models.AutoAI, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := c.col.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.AutoAI
	for cur.Next(ctx) {
		var result models.AutoAI
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

// get all AutoAI models from AutoAIsCollection by human id
func (c *AutoAIsCollectionType) GetAllByHumanID(humanID string) ([]models.AutoAI, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results []models.AutoAI
	cur, err := c.col.Find(ctx, bson.M{"humanid": humanID})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.AutoAI
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

// get a AutoAI model from AutoAIsCollection by ID
func (c *AutoAIsCollectionType) Get(id string) (models.AutoAI, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result models.AutoAI

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	err = c.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// get a AutoAI model from AutoAIsCollection by name
func (c *AutoAIsCollectionType) GetByName(name string) (models.AutoAI, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result models.AutoAI
	err := c.col.FindOne(ctx, bson.M{"name": name}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// update a AutoAI model in AutoAIsCollection
func (c *AutoAIsCollectionType) Update(data models.AutoAI) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return err
	}

	_, err = c.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})
	if err != nil {
		return err
	}
	return nil
}

// insert a AutoAI model in AutoAIsCollection
func (c *AutoAIsCollectionType) Insert(auto models.AutoAI) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := c.col.InsertOne(ctx, auto)
	if err != nil {
		return "", err
	}

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}
