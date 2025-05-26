package service

import (
	"context"

	"github.com/google/uuid"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/repository"
)

type ResetPasswordService interface {
	CreateResetPasswordRequest(
		ctx context.Context,
		userId user_entity.UserId,
	) (*entity.ResetPasswordRequest, error)
	DeleteResetPasswordRequest(
		ctx context.Context,
		validationToken uuid.UUID,
	) (*user_entity.User, error)
}

type DefaultValidationRequestService struct {
	resetPasswordRepository repository.ResetPasswordRepository
	userRepository          user_repository.UserRepository
}

func NewDefaultResetPasswordService(
	resetPasswordRepository repository.ResetPasswordRepository,
	userRepository user_repository.UserRepository,
) ResetPasswordService {
	return &DefaultValidationRequestService{
		resetPasswordRepository: resetPasswordRepository,
		userRepository:          userRepository,
	}
}

func (s *DefaultValidationRequestService) CreateResetPasswordRequest(
	ctx context.Context,
	userId user_entity.UserId,
) (*entity.ResetPasswordRequest, error) {
	validationRequestData, err := s.resetPasswordRepository.CreateResetPasswordRequest(ctx, userId)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetById(ctx, validationRequestData.UserId)
	if err != nil {
		return nil, err
	}

	return &entity.ResetPasswordRequest{
		User:  user,
		Token: validationRequestData.Token,
	}, nil
}

func (s *DefaultValidationRequestService) DeleteResetPasswordRequest(
	ctx context.Context,
	validationToken uuid.UUID,
) (*user_entity.User, error) {
	userId, err := s.resetPasswordRepository.DeleteResetPasswordRequest(ctx, validationToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
