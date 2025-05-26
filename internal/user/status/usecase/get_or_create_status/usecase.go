package get_or_create_status

import (
	"context"
	"fmt"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_status"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/save_status"
)

type GetOrCreateStatusUseCase struct {
	getStatusUseCase  *get_status.GetStatusUseCase
	saveStatusUseCase *save_status.SaveStatusUseCase
}

func NewGetOrCreateStatusUseCase(
	getStatusUseCase *get_status.GetStatusUseCase,
	saveStatusUseCase *save_status.SaveStatusUseCase,
) *GetOrCreateStatusUseCase {
	return &GetOrCreateStatusUseCase{
		getStatusUseCase:  getStatusUseCase,
		saveStatusUseCase: saveStatusUseCase,
	}
}

func (g *GetOrCreateStatusUseCase) Execute(
	ctx context.Context,
	userId user_entity.UserId,
	defaultStatus user_status_entity.Status,
) (user_status_entity.Status, error) {
	userStatus, err := g.getStatusUseCase.Execute(ctx, userId)
	if err != nil {
		newStatus := defaultStatus
		err = g.saveStatusUseCase.Execute(ctx, userId, newStatus)
		if err != nil {
			return user_status_entity.StatusUnknown, fmt.Errorf("failed to save status: %w", err)
		}
		userStatus = newStatus
	}

	return userStatus, nil
}
