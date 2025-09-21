package auction

import (
	"context"
	"fmt"
	"time"

	"github.com/Santosjordi/fcc_auction_app/configs/logger"
	"github.com/Santosjordi/fcc_auction_app/internal/entity/auction_entity"
	"github.com/Santosjordi/fcc_auction_app/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ar *AuctionRepository) FindAuctionByID(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	var auctionEntityMongo AuctioEntityMongo
	filter := bson.M{"_id": id}
	err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(fmt.Sprintf("auction with id=%s not found", id), err)
			return nil, internal_error.NewNotFoundError("auction not found", "not_found")
		}
		logger.Error("failed to find auction: ", err)
		return nil, internal_error.NewInternalServerError("failed to find auction", "internal_server_error")
	}
	auction := &auction_entity.Auction{
		ID:          auctionEntityMongo.ID,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		TimeStamp:   time.Unix(auctionEntityMongo.TimeStamp, 0),
	}
	return auction, nil
}

func (ar *AuctionRepository) FindAllAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category string,
	productName string,
) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}
	if status != 0 {
		filter["status"] = status
	}
	if category != "" {
		filter["category"] = category
	}
	if productName != "" {
		filter["product_name"] = bson.M{"$regex": productName, "$options": "i"}
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("failed to find auctions: ", err)
		return nil, internal_error.NewInternalServerError("failed to find auctions", "internal_server_error")
	}
	defer cursor.Close(ctx)

	var auctions []auction_entity.Auction
	for cursor.Next(ctx) {
		var auctionEntityMongo AuctioEntityMongo
		if err := cursor.Decode(&auctionEntityMongo); err != nil {
			logger.Error("failed to decode auction: ", err)
			return nil, internal_error.NewInternalServerError("failed to decode auction", "internal_server_error")
		}
		auction := auction_entity.Auction{
			ID:          auctionEntityMongo.ID,
			ProductName: auctionEntityMongo.ProductName,
			Category:    auctionEntityMongo.Category,
			Description: auctionEntityMongo.Description,
			Condition:   auctionEntityMongo.Condition,
			Status:      auctionEntityMongo.Status,
			TimeStamp:   time.Unix(auctionEntityMongo.TimeStamp, 0),
		}
		auctions = append(auctions, auction)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("cursor error: ", err)
		return nil, internal_error.NewInternalServerError("cursor error", "internal_server_error")
	}

	return auctions, nil
}
