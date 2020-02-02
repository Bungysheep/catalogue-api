package basemodel

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// BaseModel type
type BaseModel struct{}

// DoValidateBase - Base model validation
func (BaseModel) DoValidateBase(model interface{}) (bool, string) {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			value := modelValue.Field(i).String()

			if !isSpecified(value, field.Tag.Get("mandatory")) {
				return false, fmt.Sprintf("%s must be specified", field.Name)
			}

			if !isValidMaxLength(value, field.Tag.Get("max_length")) {
				return false, fmt.Sprintf("%s can not more than %s chars", field.Name, field.Tag.Get("max_length"))
			}

			if !isValidValue(value, field.Tag.Get("valid_value")) {
				return false, fmt.Sprintf("%s '%s' is not valid", field.Name, field.Tag.Get("valid_value"))
			}
		}
	}

	return true, ""
}

func isSpecified(value string, isMandatory string) bool {
	if isMandatory != "" {
		isMandatoryBool, _ := strconv.ParseBool(isMandatory)
		if isMandatoryBool && value == "" {
			return false
		}
	}
	return true
}

func isValidMaxLength(value string, maxLength string) bool {
	if maxLength != "" {
		maxLengthInt, _ := strconv.Atoi(maxLength)
		if len(value) > maxLengthInt {
			return false
		}
	}
	return true
}

func isValidValue(value string, validValue string) bool {
	if validValue != "" {
		values := strings.Split(validValue, ",")

		for _, item := range values {
			if item == value {
				return true
			}
		}
		return false
	}
	return true
}
