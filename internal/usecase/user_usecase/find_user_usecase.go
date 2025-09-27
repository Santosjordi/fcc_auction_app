package user_usecase

import (
	"context"

	"github.com/Santosjordi/fcc_auction_app/internal/entity/user_entity"
	"github.com/Santosjordi/fcc_auction_app/internal/internal_error"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserByID(ctx context.Context, id string) (*UserOutputDTO, internal_error.InternalError)
}

func (u *UserUseCase) FindUserByID(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError) {
	userEntity, err := u.UserRepository.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		ID:   userEntity.ID,
		Name: userEntity.Name,
	}, nil
}
