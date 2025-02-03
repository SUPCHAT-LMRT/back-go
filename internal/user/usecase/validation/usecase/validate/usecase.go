package validate

import (
	"context"
	"github.com/google/uuid"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/service"
	uberdig "go.uber.org/dig"
)

type ValidateAccountUseCaseDeps struct {
	uberdig.In
	Service        service.ValidationService
	UserRepository user_repository.UserRepository
}

type ValidateAccountUseCase struct {
	deps ValidateAccountUseCaseDeps
}

func NewValidateAccountUseCase(deps ValidateAccountUseCaseDeps) *ValidateAccountUseCase {
	return &ValidateAccountUseCase{deps: deps}
}

func (u *ValidateAccountUseCase) Execute(ctx context.Context, validationToken uuid.UUID) error {
	user, err := u.deps.Service.DeleteAccountValidationRequest(ctx, validationToken)
	if err != nil {
		return err
	}

	err = u.deps.UserRepository.SetAsVerified(ctx, user.Id)
	if err != nil {
		return err
	}

	return nil
}
