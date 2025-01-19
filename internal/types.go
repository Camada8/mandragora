package internal

// ValidationConfig holds configuration for request validation
type ValidationConfig struct {
	Body BodyConfig
}

// BodyConfig defines a map for body validation parameters
type BodyConfig map[string]string

// ValidationSet is a map of validation sets indexed by route paths
type ValidationSet map[string]Set

// ValidationShell contains the structure and parameters for validation
type ValidationShell struct {
	ValidationStruct map[string]any
	Parameters       map[string]string
}

// Set holds validation shells for body, query, and path parameters
type Set struct {
	Body      ValidationShell
	Query     ValidationShell
	Params    ValidationShell
	RoutePath string
}

// ErrorSet contains validation errors for body, query, and path parameters
type ErrorSet struct {
	BodyError   map[string]error
	QueryError  map[string]error
	ParamsError map[string]error
}
