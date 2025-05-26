package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
)

func IsUserStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	parsedStatus := entity.ParseStatus(status)

	return parsedStatus != entity.StatusUnknown
}
