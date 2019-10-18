package errors

import (
	"fmt"

	"github.com/alijared/nr-log-parser/internal/context"
)

const VALIDATION_ERROR = "ValidationError"
const EXECUTION_ERROR = "ExecutionError"

type CMDError interface {
	Error() string
	Type() string
	Usage() string
}

type ValidationError struct {
	s     string
	usage string
}

type ExecutionError struct {
	s string
}

func (e *ValidationError) Error() string {
	return e.s
}

func (e *ValidationError) Type() string {
	return VALIDATION_ERROR
}

func (e *ValidationError) Usage() string {
	return e.usage
}

func (e *ExecutionError) Error() string {
	return e.s
}

func (e *ExecutionError) Type() string {
	return EXECUTION_ERROR
}

func (e *ExecutionError) Usage() string {
	return ""
}

func NewValidationError(format string, a ...interface{}) *ValidationError {
	return &ValidationError{
		s:     fmt.Sprintf(format, a...),
		usage: context.CMD().UsageString(),
	}
}

func NewExecutionError(format string, a ...interface{}) *ExecutionError {
	return &ExecutionError{
		s: fmt.Sprintf(format, a...),
	}
}
