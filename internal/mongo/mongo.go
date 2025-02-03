package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
	"time"
)

type Client struct {
	Client *mongo.Client
}

func NewClient() (*Client, error) {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mongoClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{Client: mongoClient}, nil
}

func (r *Client) Close(ctx context.Context) error {
	return r.Client.Disconnect(ctx)
}
