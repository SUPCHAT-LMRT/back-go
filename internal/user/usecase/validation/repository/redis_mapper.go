package repository

import (
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type RedisValidationRequestMapper struct{}

func NewRedisValidationRequestMapper() mapper.Mapper[map[string]string, *ValidationRequestData] {
	return &RedisValidationRequestMapper{}
}

func (m RedisValidationRequestMapper) MapFromEntity(entity *ValidationRequestData) (map[string]string, error) {
	return map[string]string{
		"UserId": entity.UserId.String(),
		"Token":  entity.Token.String(),
	}, nil
}

func (m RedisValidationRequestMapper) MapToEntity(databaseEntity map[string]string) (*ValidationRequestData, error) {
	userId := databaseEntity["UserId"]

	validationToken, err := uuid.Parse(databaseEntity["Token"])
	if err != nil {
		return nil, err
	}

	return &ValidationRequestData{
		UserId: entity.UserId(userId),
		Token:  validationToken,
	}, nil
}
