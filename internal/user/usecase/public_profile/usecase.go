package public_profile

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type GetPublicUserProfileUseCase struct {
	repository repository.UserRepository
}

func NewGetPublicUserProfileUseCase(repository repository.UserRepository) *GetPublicUserProfileUseCase {
	return &GetPublicUserProfileUseCase{repository: repository}
}

func (u *GetPublicUserProfileUseCase) Execute(ctx context.Context, userId user_entity.UserId) (*PublicUserProfile, error) {
	user, err := u.repository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &PublicUserProfile{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

type PublicUserProfile struct {
	Id        user_entity.UserId
	FirstName string
	LastName  string
	Email     string
}
