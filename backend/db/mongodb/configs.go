package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend/models"
)

//get all OpenAIConfig configs from mongodb

func (c *ConfigsCollectionType) GetAll() ([]models.OpenAIConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := c.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find configs")
	}
	defer cur.Close(ctx)
	var configs []models.OpenAIConfig
	for cur.Next(ctx) {
		var config models.OpenAIConfig
		err := cur.Decode(&config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode config")
		}
		configs = append(configs, config)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate configs")
	}
	return configs, nil
}

// get OpenAIConfig from mongodb id
func (c *ConfigsCollectionType) GetById(id string) (models.OpenAIConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var config models.OpenAIConfig

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.OpenAIConfig{}, err
	}

	filter := bson.M{"_id": oid}
	result := c.col.FindOne(ctx, filter)

	if err := result.Decode(&config); err != nil {
		return models.OpenAIConfig{}, err
	}
	return config, nil
}

// insert new OpenAIConfig into mongodb
func (c *ConfigsCollectionType) Insert(config models.OpenAIConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.col.InsertOne(ctx, config)
	if err != nil {
		return errors.Wrap(err, "failed to insert config")
	}
	return nil
}

// update OpenAIConfig in mongodb
func (c *ConfigsCollectionType) Update(config *models.OpenAIConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(config.Id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}

	//TODO: check if the config is valid

	update := bson.M{
		"$set": bson.M{
			"name":              config.Name,
			"description":       config.Desc,
			"model":             config.Model,
			"suffix":            config.Suffix,
			"max_tokens":        config.MaxTokens,
			"temperature":       config.Temperature,
			"top_p":             config.TopP,
			"n":                 config.N,
			"stream":            config.Stream,
			"stop":              config.Stop,
			"presence_penalty":  config.PresencePenalty,
			"frequency_penalty": config.FrequencyPenalty,
			"logit_bias":        config.LogitBias,
			"user":              config.User,
		}}

	_, err = c.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to update config")
	}
	return nil
}
