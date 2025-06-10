package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func IsISO8601Date(fl validator.FieldLevel) bool {
	iso8601DateRegexString := `^[+-]?\d{4}(-(0[1-9]|1[0-2])(-(0[1-9]|[12]\d|3[01]))?)?$`
	iso8601DateRegex := regexp.MustCompile(iso8601DateRegexString)
	return iso8601DateRegex.MatchString(fl.Field().String())
}
