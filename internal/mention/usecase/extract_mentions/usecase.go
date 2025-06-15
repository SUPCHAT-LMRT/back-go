package extract_mentions

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"regexp"
)

type ExtractMentionsUseCase struct {
}

func NewExtractMentionsUseCase() *ExtractMentionsUseCase {
	return &ExtractMentionsUseCase{}
}

func (u *ExtractMentionsUseCase) Execute(messageContent string) []user_entity.UserId {
	mentionPattern := regexp.MustCompile(`<@([^>]+)>`)
	matches := mentionPattern.FindAllStringSubmatch(messageContent, -1)

	var mentions []user_entity.UserId
	for _, match := range matches {
		if len(match) > 1 {
			mentions = append(mentions, user_entity.UserId(match[1]))
		}
	}
	return mentions
}
