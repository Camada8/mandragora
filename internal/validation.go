package internal

import (
	"fmt"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

var MandragoraValidator *validator.Validate

// NewValidator initializes the MandragoraValidator with custom validation rules
func NewValidator() {
	MandragoraValidator = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validation rules for UUIDs, integers, and unsigned integers
	_ = MandragoraValidator.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil { // Return true if parsing fails (invalid UUID)
			return true
		}
		return false
	})
	_ = MandragoraValidator.RegisterValidation("integer", func(fl validator.FieldLevel) bool {
		floatNum := fl.Field().Float()
		return floatNum == math.Trunc(floatNum)
	})
	_ = MandragoraValidator.RegisterValidation("uint", func(fl validator.FieldLevel) bool {
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

// Validate performs validation on the provided parameters using the MandragoraValidator
func Validate(c *fiber.Ctx, shell ValidationShell) map[string]error {
	errors := map[string]error{}
	if MandragoraValidator == nil {
		NewValidator()
	}

	log.Infof("%+v", shell.ValidationStruct)

	for i, v := range shell.Parameters {
		log.Infof("RuleName: %s - Value: %v - Rule: %s", i, shell.ValidationStruct[i], v)
		if err := MandragoraValidator.Var(shell.ValidationStruct[i], v); err != nil {
			errors[i] = err
		}
	}
	return errors
}
