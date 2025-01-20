package mandragora

import (
	"fmt"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var mandragoraValidator *validator.Validate

// NewValidator initializes the mandragoraValidator with custom validation rules
func newValidator() {
	mandragoraValidator = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validation rules for UUIDs, integers, and unsigned integers
	_ = mandragoraValidator.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil { // Return true if parsing fails (invalid UUID)
			return true
		}
		return false
	})
	_ = mandragoraValidator.RegisterValidation("integer", func(fl validator.FieldLevel) bool {
		floatNum := fl.Field().Float()
		return floatNum == math.Trunc(floatNum)
	})
	_ = mandragoraValidator.RegisterValidation("uint", func(fl validator.FieldLevel) bool {
		floatNum := fl.Field().Float()
		return floatNum == math.Trunc(floatNum) && floatNum >= 0
	})
}

// ValidatorErrors processes validation errors and formats them into a map
func ValidatorErrors(errors map[string]error) map[string][]string {
	fields := map[string][]string{}

	if len(errors) > 0 {

		for i, err := range errors {
			err, _ := err.(validator.ValidationErrors)
			var errMsgs []string
			for _, v := range err {
				errMsgs = append(errMsgs, fmt.Sprintf("failed validation for '%s'", v.Tag()))
			}
			fields[i] = errMsgs
		}
	}

	return fields
}

// Validate performs validation on the provided parameters using the mandragoraValidator
func Validate(c *fiber.Ctx, shell ValidationShell) map[string]error {
	errors := map[string]error{}
	if mandragoraValidator == nil {
		newValidator()
	}

	for i, v := range shell.Parameters {
		if err := mandragoraValidator.Var(shell.ValidationStruct[i], v); err != nil {
			errors[i] = err
		}
	}
	return errors
}
