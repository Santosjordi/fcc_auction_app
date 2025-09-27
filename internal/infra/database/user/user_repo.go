package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Santosjordi/fcc_auction_app/configs/logger"
	"github.com/Santosjordi/fcc_auction_app/internal/entity/user_entity"
	"github.com/Santosjordi/fcc_auction_app/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: db.Collection("user"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *UserEntityMongo) error {
	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": id}
	var user UserEntityMongo
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(
				fmt.Sprintf("user with id=%s not found", id), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("user with id=%s not found", id), "not_found")
		}
		logger.Error("Error finding user by ID", err)
		return nil, internal_error.NewInternalServerError("error finding user by ID", "internal_server_error")
	}
	userEntity := &user_entity.User{
		ID:   user.Id,
		Name: user.Name,
	}
	return userEntity, nil
}
