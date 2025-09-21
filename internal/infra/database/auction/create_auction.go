package auction

import (
	"context"
	"time"

	"github.com/Santosjordi/fcc_auction_app/configs/logger"
	"github.com/Santosjordi/fcc_auction_app/internal/entity/auction_entity"
	"github.com/Santosjordi/fcc_auction_app/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctioEntityMongo struct {
	ID          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	TimeStamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auction auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctioEntityMongo{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		TimeStamp:   time.Now().Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("failed to create auction: ", err)
		return internal_error.NewInternalServerError("failed to create auction", "internal_server_error")
	}
	return nil
}
