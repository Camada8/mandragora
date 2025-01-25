package mandragora

import "reflect"

// AddValidation adds a new validation set for a specific path and kind
func AddValidation(config ValidationConfig) Set {
	return Set{
		Body:   processValidation(config.Body),
		Query:  processValidation(config.Query),
		Params: processValidation(config.Params),
	}
}

// processValidation processes the validation for a given data structure
func processValidation(data any) ValidationShell {
	t := reflect.TypeOf(data)
	parameters := make(map[string]string)
	if t == nil || t.Kind() != reflect.Struct {
		return ValidationShell{}
	}
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
