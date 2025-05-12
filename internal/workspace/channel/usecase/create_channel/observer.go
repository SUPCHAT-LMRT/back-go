package create_channel

import channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"

type CreateSpecifyChannelObserver interface {
	NotifyChannelCreated(channel *channel_entity.Channel)
}
