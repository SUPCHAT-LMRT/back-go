package utils

import (
	"fmt"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

func BuildDirectMessageRoomId(user1, user2 user_entity.UserId) string {
	// create unique room name combined to the two IDs, the room name will be the same for both users
	// so the ids are ordered
	if user1.IsAfter(user2) {
		return fmt.Sprintf("direct-%s_%s", user1.String(), user2.String())
	} else {
		return fmt.Sprintf("direct-%s_%s", user2.String(), user1.String())
	}
}
