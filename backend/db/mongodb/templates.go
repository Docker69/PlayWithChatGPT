package mongodb

import (
	"backend/models"
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//get all templates from mongodb

func (c *TemplatesCollectionType) GetAll() ([]models.Template, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := c.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find templates collection")
	}
	defer cur.Close(ctx)
	var templates []models.Template
	for cur.Next(ctx) {
		var template models.Template
		err := cur.Decode(&template)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode template document ")
		}
		templates = append(templates, template)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate templates")
	}
	return templates, nil
}

// get template by id from mongodb
func (c *TemplatesCollectionType) GetByID(id string) (models.Template, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var template models.Template
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Template{}, err
	}

	filter := bson.M{"_id": oid}
	result := c.col.FindOne(ctx, filter)
	if err := result.Decode(&template); err != nil {
		return models.Template{}, err
	}

	return template, nil
}

// get template by name from mongodb
func (c *TemplatesCollectionType) GetByName(name string) (models.Template, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var template models.Template
	filter := bson.M{"name": name}
	result := c.col.FindOne(ctx, filter)
	if err := result.Decode(&template); err != nil {
		return models.Template{}, err
	}
	return template, nil
}
