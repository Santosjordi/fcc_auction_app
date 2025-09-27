package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/Santosjordi/fcc_auction_app/configs/logger"
	"github.com/Santosjordi/fcc_auction_app/internal/entity/bid_entity"
	"github.com/Santosjordi/fcc_auction_app/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (bd *BidRepository) FindBidByAuctionID(ctx context.Context, auctionID string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionID}
	cursor, err := bd.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(
			fmt.Sprintf("failed to find bids by auction_id=%s: ", auctionID), err)
		return nil, internal_error.NewInternalServerError(
			fmt.Sprintf("failed to find bids by auction_id=%s: ", auctionID), "internal_server_error")
	}
	defer cursor.Close(ctx)

	var bids []BidEntityMongo
	if err := cursor.All(ctx, &bids); err != nil {
		logger.Error(
			fmt.Sprintf("failed to find bids by auction_id=%s: ", auctionID), err)
		return nil, internal_error.NewInternalServerError(
			fmt.Sprintf("failed to find bids by auction_id=%s: ", auctionID), "internal_server_error")
	}

	var bidEntities []bid_entity.Bid
	for _, bid := range bids {
		bidEntities = append(bidEntities, bid_entity.Bid{
			ID:        bid.ID,
			AuctionID: bid.AuctionID,
			UserID:    bid.UserID,
			Amount:    bid.Amount,
			TimeStamp: time.Unix(bid.TimeStamp, 0),
		})
	}

	return bidEntities, nil
}

func (bd *BidRepository) FindWinningBidByAuctionId(
	ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}

	var bidEntityMongo BidEntityMongo
	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})
	if err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		logger.Error("Error trying to find the auction winner", err)
		return nil, internal_error.NewInternalServerError("Error trying to find the auction winner", "internal_server_error")
	}

	return &bid_entity.Bid{
		ID:        bidEntityMongo.ID,
		UserID:    bidEntityMongo.UserID,
		AuctionID: bidEntityMongo.AuctionID,
		Amount:    bidEntityMongo.Amount,
		TimeStamp: time.Unix(bidEntityMongo.TimeStamp, 0),
	}, nil
}
