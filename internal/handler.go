package internal

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func convertMap(input map[string]string) map[string]any {
	output := make(map[string]any)
	for key, value := range input {
		output[key] = value
	}
	return output
}

func WithValidation(handler func(c *fiber.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		set := GetValidationSet(c.Route().Path)
		var errors ErrorSet
		if set.RoutePath != "" {
			if len(set.Body.Parameters) > 0 {
				c.BodyParser(&set.Body.ValidationStruct)
				log.Infof("%+v", set.Body.ValidationStruct)
				errors.BodyError = Validate(c, set.Body)
			}
			if len(set.Query.Parameters) > 0 {
				set.Query.ValidationStruct = convertMap(c.Queries())
				errors.QueryError = Validate(c, set.Query)
			}
			if len(set.Params.Parameters) > 0 {
				set.Params.ValidationStruct = convertMap(c.AllParams())
				errors.ParamsError = Validate(c, set.Params)
			}
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
		return handler(c)
	}
}
