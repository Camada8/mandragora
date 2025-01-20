package mandragora

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// convertMap converts a map of string keys and values to a map with any type values
func convertMap(input map[string]string) map[string]any {
	output := make(map[string]any)
	for key, value := range input {
		output[key] = value
	}
	return output
}

// WithValidation wraps a Fiber handler to include validation logic
func WithValidation(handler func(c *fiber.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		set := GetValidationSet(c.Route().Path)
		var errors ErrorSet
		if set.RoutePath != "" {
			// Validate request body if parameters are defined
			if len(set.Body.Parameters) > 0 {
				c.BodyParser(&set.Body.ValidationStruct)
				errors.BodyError = Validate(c, set.Body)
			}
			// Validate query parameters
			if len(set.Query.Parameters) > 0 {
				set.Query.ValidationStruct = convertMap(c.Queries())
				errors.QueryError = Validate(c, set.Query)
			}
			// Validate path parameters
			if len(set.Params.Parameters) > 0 {
				set.Params.ValidationStruct = convertMap(c.AllParams())
				errors.ParamsError = Validate(c, set.Params)
			}
			// If there are validation errors, respond with a bad request
			if len(errors.BodyError) > 0 || len(errors.QueryError) > 0 || len(errors.ParamsError) > 0 {
				c.Status(fiber.StatusBadRequest)
				c.Type("json", "utf-8")
				j, _ := json.Marshal(fiber.Map{
					"error": true,
					"msg": fiber.Map{
						"body":   ValidatorErrors(errors.BodyError),
						"query":  ValidatorErrors(errors.QueryError),
						"params": ValidatorErrors(errors.ParamsError),
					},
				})
				return c.Send(j)
			}
		}
		return handler(c) // Call the original handler if validation passes
	}
}
