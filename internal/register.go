package internal

import "reflect"

var ValidationSets = make(ValidationSet)

func GetValidationSet(path string) Set {
	return ValidationSets[path]
}

func GetValidationSets() ValidationSet {
	return ValidationSets
}

func AddValidationSet(path string, data any, kind string) {
	if kind == "body" || kind == "query" || kind == "params" {
		t := reflect.TypeOf(data)
		parameters := make(map[string]string)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// Get the struct tags
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

func SetBodyValidation(path string, validator any) {
	AddValidationSet(path, validator, "body")
}

func SetQueryValidation(path string, validator any) {
	AddValidationSet(path, validator, "query")
}

func SetParamsValidation(path string, validator any) {
	AddValidationSet(path, validator, "params")
}
