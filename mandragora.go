package mandragora

import "github.com/ErnestoMuniz/mandragora/internal"

// Expose validation functions from the internal package
var GetValidationSets = internal.GetValidationSets
var SetBodyValidation = internal.SetBodyValidation
var SetQueryValidation = internal.SetQueryValidation
var SetParamsValidation = internal.SetParamsValidation
var WithValidation = internal.WithValidation
