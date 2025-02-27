package util

import (
	"fmt"
	"strings"

	valid "github.com/go-playground/validator/v10"
	"github.com/hewpao/hewpao-backend/types"
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

func CategoryValidator(fl valid.FieldLevel) bool {
	category := fl.Field().String()
	switch category {
	case string(types.Electronics), string(types.Fashion), string(types.Food), string(types.Books), string(types.Other):
		return true
	}
	return false
}

func DeliveryStatusValidator(fl valid.FieldLevel) bool {
	deliveryStatyus := fl.Field().String()
	switch deliveryStatyus {
	case string(types.Pending), string(types.Purchased), string(types.PickedUp), string(types.OutForDelivery), string(types.Cancel), string(types.Returned), string(types.Refunded), string(types.Opening), string(types.Delivered):
		return true
	}
	return false
}

func ValidateStruct[T any](payload T) *ValidateError {
	validate.RegisterValidation("category", CategoryValidator)
	validate.RegisterValidation("delivery-status", DeliveryStatusValidator)

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
