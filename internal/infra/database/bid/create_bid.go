package bid

import (
	"context"
	"sync"

	"github.com/Santosjordi/fcc_auction_app/configs/logger"
	"github.com/Santosjordi/fcc_auction_app/internal/entity/auction_entity"
	"github.com/Santosjordi/fcc_auction_app/internal/entity/bid_entity"
	"github.com/Santosjordi/fcc_auction_app/internal/infra/database/auction"
	"github.com/Santosjordi/fcc_auction_app/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	ID        string  `bson:"_id"`
	AuctionID string  `bson:"auction_id"`
	UserID    string  `bson:"user_id"`
	Amount    float64 `bson:"amount"`
	TimeStamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *auction.AuctionRepository
}

func (bd *BidRepository) CreateBid(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, bid := range bidEntities {
		wg.Add(1)

		go func(bidValue bid_entity.Bid) {
			defer wg.Done()

			auctionEntity, err := bd.AuctionRepository.FindAuctionByID(ctx, bidValue.AuctionID)
			if err != nil {
				logger.Error("Error trying to find auction by id", err)
				return
			}

			if auctionEntity.Status != auction_entity.Active {
				logger.Error("Auction is not active", nil)
				return
			}

			bidEntityMongo := &BidEntityMongo{
				ID:        bidValue.ID,
				AuctionID: bidValue.AuctionID,
				UserID:    bidValue.UserID,
				Amount:    bidValue.Amount,
				TimeStamp: bidValue.TimeStamp.Unix(),
			}

			if _, err := bd.Collection.InsertOne(ctx, bidEntityMongo); err != nil {
				logger.Error("failed to create bid: ", err)
				return
			}
		}(bid)
	}

	wg.Wait()
	return nil
}
