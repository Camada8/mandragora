package mandragora

import "reflect"

// Global variable to hold validation sets
var ValidationSets = make(ValidationSet)

// GetValidationSet retrieves a validation set for a given path
func GetValidationSet(path string) Set {
	return ValidationSets[path]
}

// GetValidationSets returns all validation sets
func GetValidationSets() ValidationSet {
	return ValidationSets
}

// AddValidation adds a new validation set for a specific path and kind
func AddValidation(config ValidationConfig) {
	set := ValidationSets[config.RoutePath]
	set.RoutePath = config.RoutePath

	// Process validation for body, query, and params
	set.Body = processValidation(config.Body)
	set.Query = processValidation(config.Query)
	set.Params = processValidation(config.Params)

	ValidationSets[config.RoutePath] = set
}

// processValidation processes the validation for a given data structure
func processValidation(data interface{}) ValidationShell {
	t := reflect.TypeOf(data)
	parameters := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Extract struct tags for query, json, and validation
		queryTag := field.Tag.Get("query")
		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("validate")
		if queryTag != "" {
			parameters[queryTag] = validateTag
		} else if jsonTag != "" {
			parameters[jsonTag] = validateTag
		} else {
			parameters[field.Name] = validateTag
		}
	}
	typedData, _ := data.(map[string]any)
	return ValidationShell{
		ValidationStruct: typedData,
		Parameters:       parameters,
	}
}
