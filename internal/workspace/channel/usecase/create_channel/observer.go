package create_channel

import channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"

type CreateChannelObserver interface {
	ChannelCreated(channel *channel_entity.Channel)
}
