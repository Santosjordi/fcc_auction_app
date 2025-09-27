package bid_entity

import "time"

type Bid struct {
	ID        string
	AuctionID string
	UserID    string
	Amount    float64
	TimeStamp time.Time
}
