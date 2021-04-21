package helpers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"log"
)

// Validate the value of toValidate using
// be sure to pass a reference in the toValidate field
// second return value is true if there are no errors
func Validate(c echo.Context, toValidate interface{}) (map[string]string, bool) {
	err := c.Validate(toValidate)

	// validationData contains a map of field names to error strings
	validationData := map[string]string{}
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Printf("invalidValidationError - %s\n", err)

		} else {
			var valErrs validator.ValidationErrors
			if errors.As(err, &valErrs) {
				for _, err := range valErrs {
					errString := ""
					if err.Tag() == "required" {
						errString = fmt.Sprintf("%s is required", err.Field())
					} else {
						errString = fmt.Sprintf("'%s' failed the '%s' check", err.Field(), err.Tag())
					}

					validationData[err.Field()] = errString
				}
				SetErrorFlash(c, "Form data is invalid. Please correct the errors below.")
			} else {
				log.Printf("not a validator.ValidationErrors")
			}
		}
	}
	return validationData, len(validationData) == 0
}
