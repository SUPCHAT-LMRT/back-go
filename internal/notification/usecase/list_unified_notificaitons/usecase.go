package list_unified_notifications

//
//import (
//	"context"
//	"fmt"
//	"github.com/supchat-lmrt/back-go/internal/notification/entity"
//	"github.com/supchat-lmrt/back-go/internal/notification/repository"
//	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
//	"go.uber.org/dig"
//	"sort"
//	"strings"
//)
//
//type ListUnifiedNotificationsUseCaseDeps struct {
//	dig.In
//	DirectMessageRepo   repository.DirectMessageNotificationRepository
//	ChannelMessageRepo  repository.ChannelMessageNotificationRepository
//	WorkspaceInviteRepo repository.WorkspaceInviteNotificationRepository
//}
//
//type ListUnifiedNotificationsUseCase struct {
//	deps ListUnifiedNotificationsUseCaseDeps
//}
//
//func NewListUnifiedNotificationsUseCase(deps ListUnifiedNotificationsUseCaseDeps) *ListUnifiedNotificationsUseCase {
//	return &ListUnifiedNotificationsUseCase{deps: deps}
//}
//
//type ListUnifiedNotificationsRequest struct {
//	UserId user_entity.UserId
//	Limit  int
//	Offset int
//}
//
//type ListUnifiedNotificationsResponse struct {
//	Notifications []*entity.UnifiedNotification `json:"notifications"`
//	Total         int                           `json:"total"`
//	UnreadCount   int                           `json:"unread_count"`
//}
//
//func (u *ListUnifiedNotificationsUseCase) Execute(ctx context.Context, req ListUnifiedNotificationsRequest) (*ListUnifiedNotificationsResponse, error) {
//	// Récupérer toutes les notifications en parallèle
//	directNotifications, err := u.deps.DirectMessageRepo.List(ctx, req.UserId)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get direct message notifications: %w", err)
//	}
//
//	channelNotifications, err := u.deps.ChannelMessageRepo.List(ctx, req.UserId)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get channel message notifications: %w", err)
//	}
//
//	workspaceInviteNotifications, err := u.deps.WorkspaceInviteRepo.List(ctx, req.UserId)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get workspace invite notifications: %w", err)
//	}
//
//	// Convertir toutes les notifications en format unifié
//	var unifiedNotifications []*entity.UnifiedNotification
//
//	// Convertir les notifications de messages directs
//	for _, notif := range directNotifications {
//		unified := u.mapDirectMessageToUnified(notif)
//		unifiedNotifications = append(unifiedNotifications, unified)
//	}
//
//	// Convertir les notifications de messages de canaux
//	for _, notif := range channelNotifications {
//		unified := u.mapChannelMessageToUnified(notif)
//		unifiedNotifications = append(unifiedNotifications, unified)
//	}
//
//	// Convertir les notifications d'invitations de workspace
//	for _, notif := range workspaceInviteNotifications {
//		unified := u.mapWorkspaceInviteToUnified(notif)
//		unifiedNotifications = append(unifiedNotifications, unified)
//	}
//
//	// Trier par date de création (plus récent en premier)
//	sort.Slice(unifiedNotifications, func(i, j int) bool {
//		return unifiedNotifications[i].CreatedAt.After(unifiedNotifications[j].CreatedAt)
//	})
//
//	// Compter les non-lues
//	unreadCount := 0
//	for _, notif := range unifiedNotifications {
//		if !notif.IsRead {
//			unreadCount++
//		}
//	}
//
//	// Appliquer la pagination
//	total := len(unifiedNotifications)
//	if req.Offset >= total {
//		unifiedNotifications = []*entity.UnifiedNotification{}
//	} else {
//		end := req.Offset + req.Limit
//		if end > total {
//			end = total
//		}
//		unifiedNotifications = unifiedNotifications[req.Offset:end]
//	}
//
//	return &ListUnifiedNotificationsResponse{
//		Notifications: unifiedNotifications,
//		Total:         total,
//		UnreadCount:   unreadCount,
//	}, nil
//}
//
//func (u *ListUnifiedNotificationsUseCase) mapDirectMessageToUnified(notif *entity.DirectMessageNotification) *entity.UnifiedNotification {
//	return &entity.UnifiedNotification{
//		Id:          notif.Id,
//		Type:        entity.NotificationTypeDirectMessage,
//		Title:       fmt.Sprintf("Message de %s", notif.SenderName),
//		Description: u.truncateMessage(notif.MessagePreview),
//		ImageUrl:    "", // ajouter l'avatar de l'utilisateur
//		IsRead:      notif.IsRead,
//		CreatedAt:   notif.CreatedAt,
//		DirectMessageData: &entity.DirectMessageNotificationData{
//			SenderId:       notif.SenderId,
//			SenderName:     notif.SenderName,
//			MessageId:      notif.MessageId,
//			MessagePreview: notif.MessagePreview,
//		},
//	}
//}
//
//func (u *ListUnifiedNotificationsUseCase) mapChannelMessageToUnified(notif *entity.ChannelMessageNotification) *entity.UnifiedNotification {
//	return &entity.UnifiedNotification{
//		Id:          notif.Id,
//		Type:        entity.NotificationTypeChannelMessage,
//		Title:       fmt.Sprintf("%s dans #%s", notif.SenderName, notif.ChannelName),
//		Description: u.truncateMessage(notif.MessagePreview),
//		ImageUrl:    "", // ajouter l'avatar de l'utilisateur ou l'icône du workspace
//		IsRead:      notif.IsRead,
//		CreatedAt:   notif.CreatedAt,
//		ChannelMessageData: &entity.ChannelMessageNotificationData{
//			SenderId:       notif.SenderId,
//			SenderName:     notif.SenderName,
//			ChannelId:      notif.ChannelId,
//			ChannelName:    notif.ChannelName,
//			WorkspaceId:    notif.WorkspaceId,
//			WorkspaceName:  notif.WorkspaceName,
//			MessageId:      notif.MessageId,
//			MessagePreview: notif.MessagePreview,
//		},
//	}
//}
//
//func (u *ListUnifiedNotificationsUseCase) mapWorkspaceInviteToUnified(notif *entity.WorkspaceInviteNotification) *entity.UnifiedNotification {
//	return &entity.UnifiedNotification{
//		Id:          notif.Id,
//		Type:        entity.NotificationTypeWorkspaceInvite,
//		Title:       "Invitation à un workspace",
//		Description: fmt.Sprintf("%s vous a invité à rejoindre %s", notif.InviterName, notif.WorkspaceName),
//		ImageUrl:    "", // ajouter l'icône du workspace
//		IsRead:      notif.IsRead,
//		CreatedAt:   notif.CreatedAt,
//		WorkspaceInviteData: &entity.WorkspaceInviteNotificationData{
//			InviterId:     notif.InviterId,
//			InviterName:   notif.InviterName,
//			WorkspaceId:   notif.WorkspaceId,
//			WorkspaceName: notif.WorkspaceName,
//		},
//	}
//}
//
//func (u *ListUnifiedNotificationsUseCase) truncateMessage(message string, maxLength ...int) string {
//	length := 100
//	if len(maxLength) > 0 {
//		length = maxLength[0]
//	}
//
//	if len(message) <= length {
//		return message
//	}
//
//	// Trouver le dernier espace avant la limite pour éviter de couper un mot
//	truncated := message[:length]
//	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > length/2 {
//		truncated = truncated[:lastSpace]
//	}
//
//	return truncated + "..."
//}
