package internal

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

// AddValidationSet adds a new validation set for a specific path and kind
func AddValidationSet(path string, data any, kind string) {
	if kind == "body" || kind == "query" || kind == "params" {
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
		set := ValidationSets[path]
		set.RoutePath = path
		typedData, _ := data.(map[string]any)
		// Assign validation shell based on the kind
		if kind == "body" {
			set.Body = ValidationShell{
				ValidationStruct: typedData,
				Parameters:       parameters,
			}
		}
		if kind == "query" {
			set.Query = ValidationShell{
				ValidationStruct: typedData,
				Parameters:       parameters,
			}
		}
		if kind == "params" {
			set.Params = ValidationShell{
				ValidationStruct: typedData,
				Parameters:       parameters,
			}
		}
		ValidationSets[path] = set
	}
}

// SetBodyValidation sets validation for the body of a request
func SetBodyValidation(path string, validator any) {
	AddValidationSet(path, validator, "body")
}

// SetQueryValidation sets validation for the query parameters of a request
func SetQueryValidation(path string, validator any) {
	AddValidationSet(path, validator, "query")
}

// SetParamsValidation sets validation for the path parameters of a request
func SetParamsValidation(path string, validator any) {
	AddValidationSet(path, validator, "params")
}
