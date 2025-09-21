package mongo

import (
	"context"
	"os"

	"github.com/Santosjordi/fcc_auction_app/configs/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_URL = "MONGODB_URL"
	MONGODB_DB  = "MONGODB_DB"
)

func NewMongoDBConection(ctx context.Context) (*mongo.Database, error) {
	mongoURL := os.Getenv(MONGODB_URL)
	mongoDB := os.Getenv(MONGODB_DB)

	client, err := mongo.Connect(
		ctx, options.Client().ApplyURI(mongoURL),
	)

	if err != nil {
		logger.Error("Error connecting to MongoDB", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Error pinging MongoDB", err)
		return nil, err
	}

	return client.Database(mongoDB), nil

}
