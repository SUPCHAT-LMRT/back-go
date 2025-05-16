package update_user

import (
	"context"
	"fmt"
	user_search "github.com/supchat-lmrt/back-go/internal/search/user"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type UpdateUserUseCase struct {
	repository            repository.UserRepository
	searchUserSyncManager user_search.SearchUserSyncManager
}

func NewUpdateUserUseCase(repository repository.UserRepository, searchUserSyncManager user_search.SearchUserSyncManager) *UpdateUserUseCase {
	return &UpdateUserUseCase{repository: repository, searchUserSyncManager: searchUserSyncManager}
}

func (u *UpdateUserUseCase) Execute(ctx context.Context, user *user_entity.User) error {
	err := u.repository.Update(ctx, user)
	if err != nil {
		return err
	}

	err = u.searchUserSyncManager.AddUser(ctx, &user_search.SearchUser{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("error syncing user: %w", err)
	}

	return err
}

func (u *UpdateUserUseCase) GetUserById(ctx context.Context, userId user_entity.UserId) (*user_entity.User, error) {
	return u.repository.GetById(ctx, userId)
}
