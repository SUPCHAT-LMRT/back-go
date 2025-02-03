package service

import (
	"context"
	"github.com/google/uuid"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/repository"
)

type ValidationService interface {
	CreateAccountValidationRequest(ctx context.Context, userId user_entity.UserId) (*entity.ValidationRequest, error)
	DeleteAccountValidationRequest(ctx context.Context, validationToken uuid.UUID) (*user_entity.User, error)
}

type DefaultValidationRequestService struct {
	validationRequestRepository repository.ValidationRepository
	userRepository              user_repository.UserRepository
}

func NewDefaultValidationRequestService(validationRequestRepository repository.ValidationRepository, userRepository user_repository.UserRepository) ValidationService {
	return &DefaultValidationRequestService{
		validationRequestRepository: validationRequestRepository,
		userRepository:              userRepository,
	}
}

func (s *DefaultValidationRequestService) CreateAccountValidationRequest(ctx context.Context, userId user_entity.UserId) (*entity.ValidationRequest, error) {
	validationRequestData, err := s.validationRequestRepository.CreateValidationRequest(ctx, userId)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetById(ctx, validationRequestData.UserId)
	if err != nil {
		return nil, err
	}

	return &entity.ValidationRequest{
		User:  user,
		Token: validationRequestData.Token,
	}, nil
}

func (s *DefaultValidationRequestService) DeleteAccountValidationRequest(ctx context.Context, validationToken uuid.UUID) (*user_entity.User, error) {
	userId, err := s.validationRequestRepository.DeleteValidationRequest(ctx, validationToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
