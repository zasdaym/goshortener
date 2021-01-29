package link

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "links"

// MongoService is Link Service with mongoDB implementation.
type MongoService struct {
	db *mongo.Database
}

// NewMongoService creates new MongoService.
func NewMongoService(db *mongo.Database) *MongoService {
	return &MongoService{db: db}
}

func (s *MongoService) getLinks(ctx context.Context, req getLinksRequest) (links []*link, err error) {
	links = make([]*link, 0)
	skip := req.Limit * (req.Page - 1)
	opts := options.FindOptions{
		Limit: &req.Limit,
		Skip:  &skip,
	}
	cur, err := s.db.Collection(collectionName).Find(ctx, bson.M{}, &opts)
	if err != nil {
		return links, fmt.Errorf("failed to get links: %w", err)
	}
	for cur.Next(ctx) {
		var l link
		if err := cur.Decode(&l); err != nil {
			return links, fmt.Errorf("failed to decode link: %w", err)
		}
		links = append(links, &l)
	}
	return links, err
}

func (s *MongoService) createLink(ctx context.Context, req createLinkRequest) (err error) {
	l := link{
		ID:        primitive.NewObjectID(),
		Slug:      req.Slug,
		Target:    req.Target,
		CreatedAt: time.Now(),
	}
	if _, err := s.db.Collection(collectionName).InsertOne(ctx, l); err != nil {
		return fmt.Errorf("failed to insert link to db: %w", err)
	}
	return nil
}
