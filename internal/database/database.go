package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"beanstock/internal/types"
)

type Service interface {
	Health() map[string]string
	StoreWebsite(site types.Website) error
	GetWebsite(id string) (types.Website, error)
}

type service struct {
	db *mongo.Client
}

var (
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	//database = os.Getenv("DB_DATABASE")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db: client,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) StoreWebsite(site types.Website) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	coll := s.db.Database("beanstock").Collection("website")
	_, err := coll.InsertOne(ctx, site)

	return err
}

func (s *service) GetWebsite(id string) (site types.Website, e error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	coll := s.db.Database("beanstock").Collection("website")
	filter := bson.D{{Key: "id", Value: id}}
	e = coll.FindOne(ctx, filter).Decode(&site)

	return site, e
}
