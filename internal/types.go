package internal

type ValidationConfig struct {
	Body BodyConfig
}

type BodyConfig map[string]string

type ValidationSet map[string]Set

type ValidationShell struct {
	ValidationStruct map[string]any
	Parameters       map[string]string
}

type Set struct {
	Body      ValidationShell
	Query     ValidationShell
	Params    ValidationShell
	RoutePath string
}

type ErrorSet struct {
	BodyError   map[string]error
	QueryError  map[string]error
	ParamsError map[string]error
}
