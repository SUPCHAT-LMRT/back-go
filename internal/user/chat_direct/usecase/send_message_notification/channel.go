package send_message_notification

import "context"

type Channel interface {
	SendNotification(ctx context.Context, req SendMessageNotificationRequest) error
}
