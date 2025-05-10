package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type OauthConnectionId string

type OauthConnection struct {
	Id          OauthConnectionId
	UserId      user_entity.UserId
	Provider    string
	OauthEmail  string
	OauthUserId string
}

func (id OauthConnectionId) String() string {
	return string(id)
}
