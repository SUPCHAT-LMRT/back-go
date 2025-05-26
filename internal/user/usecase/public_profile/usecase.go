package public_profile

import (
	"context"
	"fmt"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
	uberdig "go.uber.org/dig"
)

type GetPublicUserProfileUseCaseDeps struct {
	uberdig.In
	UserRepository         repository.UserRepository
	GetPublicStatusUseCase *get_public_status.GetPublicStatusUseCase
}

type GetPublicUserProfileUseCase struct {
	deps GetPublicUserProfileUseCaseDeps
}

func NewGetPublicUserProfileUseCase(
	deps GetPublicUserProfileUseCaseDeps,
) *GetPublicUserProfileUseCase {
	return &GetPublicUserProfileUseCase{deps: deps}
}

func (u GetPublicUserProfileUseCase) Execute(
	ctx context.Context,
	userId user_entity.UserId,
) (*PublicUserProfile, error) {
	user, err := u.deps.UserRepository.GetById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	status, err := u.deps.GetPublicStatusUseCase.Execute(ctx, userId, entity.StatusOffline)
	if err != nil {
		return nil, fmt.Errorf("unable to get user status: %w", err)
	}

	return &PublicUserProfile{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Status:    status,
	}, nil
}

type PublicUserProfile struct {
	Id        user_entity.UserId
	FirstName string
	LastName  string
	Email     string
	Status    entity.Status
}
