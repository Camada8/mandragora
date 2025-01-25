package mandragora

// ValidationConfig holds configuration for request validation
type ValidationConfig struct {
	Body   any
	Query  any
	Params any
}

// ValidationSet is a map of validation sets indexed by route paths
type ValidationSet map[string]Set

// ValidationShell contains the structure and parameters for validation
type ValidationShell struct {
	ValidationStruct map[string]any
	Parameters       map[string]string
}

// Set holds validation shells for body, query, and path parameters
type Set struct {
	Body   ValidationShell
	Query  ValidationShell
	Params ValidationShell
}

// ErrorSet contains validation errors for body, query, and path parameters
type ErrorSet struct {
	BodyError   map[string]error
	QueryError  map[string]error
	ParamsError map[string]error
}
