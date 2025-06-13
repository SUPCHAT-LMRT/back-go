package export_user_data

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	"time"
)

type ExportUserDataUseCase struct {
	userRepository repository.UserRepository
}

func NewExportUserDataUseCase(
	userRepository repository.UserRepository,
) *ExportUserDataUseCase {
	return &ExportUserDataUseCase{userRepository: userRepository}
}

func (u *ExportUserDataUseCase) Execute(ctx context.Context, userID entity.UserId) (ExportableUserData, error) {
	user, err := u.userRepository.GetById(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ExportableUserData{}, repository.ErrUserNotFound
		}
		return ExportableUserData{}, err
	}

	data := ExportableUserData{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return data, nil
}

type ExportableUserData struct {
	Id        entity.UserId `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
