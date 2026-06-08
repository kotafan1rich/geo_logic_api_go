package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ParseValidationError(err error) string {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		var errorMessages []string
		for _, fe := range ve {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", fe.Field(), fe.Tag()))
		}
		return "Field validation failed on fields: " + strings.Join(errorMessages, ", ")

	}
	out := "Malformed JSON structure or invalid data types"
	return out
}
