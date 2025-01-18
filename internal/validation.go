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

// NewValidator func for create a new validator for model fields.
func NewValidator() {
	// Create a new validator for a Book model.
	MandragoraValidator = validator.New(validator.WithRequiredStructEnabled())

	// Custom validation for uuid.UUID fields.
	_ = MandragoraValidator.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
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

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(errors map[string]error) map[string][]string {
	// Define fields map.
	fields := map[string][]string{}

	// Make error message for each invalid field.
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
