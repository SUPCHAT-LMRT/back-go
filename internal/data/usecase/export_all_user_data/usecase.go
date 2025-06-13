package export_all_user_data

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/export_data_chat_direct"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/export_user_data"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/export_data_chat_message"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/export_data_user_workspace"
	uberdig "go.uber.org/dig"
)

type ExportAllUserDataUseCaseDeps struct {
	uberdig.In
	ExportUserDataUseCase          *export_user_data.ExportUserDataUseCase
	ExportDataChatDirectUseCase    *export_data_chat_direct.ExportDataChatDirectUseCase
	ExportDataChatMessageUseCase   *export_data_chat_message.ExportDataChatMessageUseCase
	ExportDataUserWorkspaceUseCase *export_data_user_workspace.ExportDataUserWorkspaceUseCase
}

type ExportAllUserDataUseCase struct {
	deps ExportAllUserDataUseCaseDeps
}

func NewExportAllUserDataUseCase(
	deps ExportAllUserDataUseCaseDeps,
) *ExportAllUserDataUseCase {
	return &ExportAllUserDataUseCase{deps: deps}
}

func (u *ExportAllUserDataUseCase) Execute(ctx context.Context, user entity.UserId) ([]byte, error) {
	userProfileData, err := u.deps.ExportUserDataUseCase.Execute(ctx, user)
	if err != nil {
		return nil, err
	}

	directMessages, err := u.deps.ExportDataChatDirectUseCase.Execute(ctx, user)
	if err != nil {
		return nil, err
	}

	chatMessages, err := u.deps.ExportDataChatMessageUseCase.Execute(ctx, user)
	if err != nil {
		return nil, err
	}

	userWorkspaces, err := u.deps.ExportDataUserWorkspaceUseCase.Execute(ctx, user)
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"userProfileData": userProfileData,
		"directMessages":  directMessages,
		"chatMessages":    chatMessages,
		"userWorkspaces":  userWorkspaces,
	}

	exportedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	return exportedData, nil
}
