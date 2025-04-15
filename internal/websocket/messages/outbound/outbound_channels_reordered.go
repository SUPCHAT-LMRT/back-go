package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type OutboundChannelsReordered struct {
	messages.DefaultMessage
	ChannelReorders []ChannelReorderMessage `json:"channelReorders"`
}

type ChannelReorderMessage struct {
	ChannelId entity.ChannelId `json:"channelId"`
	NewOrder  int              `json:"newOrder"`
}

func (m OutboundChannelsReordered) GetActionName() messages.Action {
	return messages.OutboundChannelsReorderedAction
}

func (m OutboundChannelsReordered) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
