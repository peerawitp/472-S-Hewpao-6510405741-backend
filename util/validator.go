package util

import (
	"fmt"
	"strings"

	valid "github.com/go-playground/validator/v10"
)

var validate = valid.New()

type ValidateError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

var customErrorMessages = map[string]string{
	"min": "must be at least %s",
	"max": "must be at most %s",
}

func getCustomErrorMessage(tag string, param string) string {
	if msg, ok := customErrorMessages[tag]; ok {
		return fmt.Sprintf(msg, param)
	}
	return ""
}

func ValidateStruct[T any](payload T) *ValidateError {
	err := validate.Struct(payload)
	errMsg := ""
	if err != nil {
		for _, err := range err.(valid.ValidationErrors) {
			field := strings.Split(err.StructNamespace(), ".")[len(strings.Split(err.StructNamespace(), "."))-1]
			tag := err.Tag()
			param := err.Param()

			var msg string
			if customMsg := getCustomErrorMessage(tag, param); customMsg != "" {
				msg = fmt.Sprintf("%s %s", field, customMsg)
			} else {
				msg = fmt.Sprintf("%s is %s", field, tag)
			}

			msg = strings.ToLower(string(msg[0])) + msg[1:]
			errMsg = errMsg + msg + ", "
		}

		return &ValidateError{
			Error:   "invalid request",
			Message: errMsg[:len(errMsg)-2],
		}
	}

	return nil
}
