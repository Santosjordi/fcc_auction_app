package user_entity

import (
	"context"

	"github.com/Santosjordi/fcc_auction_app/internal/internal_error"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserRepositoryInterface interface {
	FindUserByID(ctx context.Context, id string) (*User, *internal_error.InternalError)
}
