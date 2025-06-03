package service

import (
	"context"

	"github.com/google/uuid"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/repository"
)

type ForgotPasswordService interface {
	CreateForgotPasswordRequest(
		ctx context.Context,
		userId user_entity.UserId,
	) (*entity.ForgotPasswordRequest, error)
	DeleteForgotPasswordRequest(
		ctx context.Context,
		validationToken uuid.UUID,
	) (*user_entity.User, error)
}

type DefaultValidationRequestService struct {
	forgotPasswordRequestRepository repository.ForgotPasswordRepository
	userRepository                  user_repository.UserRepository
}

func NewDefaultForgotPasswordRequestService(
	validationRequestRepository repository.ForgotPasswordRepository,
	userRepository user_repository.UserRepository,
) ForgotPasswordService {
	return &DefaultValidationRequestService{
		forgotPasswordRequestRepository: validationRequestRepository,
		userRepository:                  userRepository,
	}
}

func (s *DefaultValidationRequestService) CreateForgotPasswordRequest(
	ctx context.Context,
	userId user_entity.UserId,
) (*entity.ForgotPasswordRequest, error) {
	validationRequestData, err := s.forgotPasswordRequestRepository.CreateForgotPasswordRequest(
		ctx,
		userId,
	)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetById(ctx, validationRequestData.UserId)
	if err != nil {
		return nil, err
	}

	return &entity.ForgotPasswordRequest{
		User:  user,
		Token: validationRequestData.Token,
	}, nil
}

func (s *DefaultValidationRequestService) DeleteForgotPasswordRequest(
	ctx context.Context,
	validationToken uuid.UUID,
) (*user_entity.User, error) {
	userId, err := s.forgotPasswordRequestRepository.DeleteForgotPasswordRequest(
		ctx,
		validationToken,
	)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
