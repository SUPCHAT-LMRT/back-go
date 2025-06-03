package mark_notification_as_read

//
//import (
//	"github.com/supchat-lmrt/back-go/internal/notification/entity"
//	"github.com/supchat-lmrt/back-go/internal/notification/repository"
//	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
//	"go.uber.org/dig"
//)
//
//type MarkNotificationAsReadUseCaseDeps struct {
//	dig.In
//	DirectMessageRepo   repository.DirectMessageNotificationRepository
//	ChannelMessageRepo  repository.ChannelMessageNotificationRepository
//	WorkspaceInviteRepo repository.WorkspaceInviteNotificationRepository
//}
//
//type MarkNotificationAsReadUseCase struct {
//	deps MarkNotificationAsReadUseCaseDeps
//}
//
//func NewMarkNotificationAsReadUseCase(deps MarkNotificationAsReadUseCaseDeps) *MarkNotificationAsReadUseCase {
//	return &MarkNotificationAsReadUseCase{deps: deps}
//}
//
//type MarkNotificationAsReadRequest struct {
//	UserId         user_entity.UserId
//	NotificationId entity.Notification
