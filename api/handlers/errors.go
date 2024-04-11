package handlers

import "fmt"

// Not valid Id Error

type notValidParamError struct {
	Param string
	Err   error
}

func (e *notValidParamError) Error() string {
	return fmt.Sprintf("Not a valid param '%s': \n%v", e.Param, e.Err.Error())
}

func (e *notValidParamError) Msg() string {
	return fmt.Sprintf("Not a valid param '%s'", e.Param)
}

// Operation Error

type operationError struct {
	Entity    string
	Operation operation
	origin    error
}

func (e *operationError) Error() string {
	return fmt.Sprintf("Something went wrong %s the %s: %v", e.Operation, e.Entity, e.origin.Error())
}

func (e *operationError) Msg() string {
	return fmt.Sprintf("Something went wrong %s the %s.", e.Operation, e.Entity)
}

// JSON binding Error

type jsonBindingError struct {
	Err error
}

func (e *jsonBindingError) Error() string {
	return fmt.Sprintf("Invalid request format. Please provide valid JSON data: \n%v", e.Err.Error())
}

func (e *jsonBindingError) Msg() string {
	return "Invalid request format. Please provide valid JSON data."
}
