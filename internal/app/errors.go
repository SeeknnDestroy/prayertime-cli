package app

import "errors"

const (
	ExitSuccess  = 0
	ExitFailure  = 1
	ExitUsage    = 2
	ExitNotFound = 3
	ExitNetwork  = 4
	ExitConflict = 5
)

type CLIError struct {
	ExitCode      int
	ErrorType     string
	Message       string
	InputReceived string
	Suggestion    string
	Cause         error
}

func (e *CLIError) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func (e *CLIError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func AsCLIError(err error) *CLIError {
	if err == nil {
		return nil
	}

	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr
	}

	return &CLIError{
		ExitCode:  ExitFailure,
		ErrorType: "internal_error",
		Message:   err.Error(),
		Cause:     err,
	}
}

func NewUsageError(message, input, suggestion string) *CLIError {
	return &CLIError{
		ExitCode:      ExitUsage,
		ErrorType:     "usage_error",
		Message:       message,
		InputReceived: input,
		Suggestion:    suggestion,
	}
}

func NewNotFoundError(message, input, suggestion string) *CLIError {
	return &CLIError{
		ExitCode:      ExitNotFound,
		ErrorType:     "not_found",
		Message:       message,
		InputReceived: input,
		Suggestion:    suggestion,
	}
}

func NewAmbiguousError(message, input, suggestion string) *CLIError {
	return &CLIError{
		ExitCode:      ExitNotFound,
		ErrorType:     "ambiguous_location",
		Message:       message,
		InputReceived: input,
		Suggestion:    suggestion,
	}
}

func NewNetworkError(message, input, suggestion string, cause error) *CLIError {
	return &CLIError{
		ExitCode:      ExitNetwork,
		ErrorType:     "network_timeout",
		Message:       message,
		InputReceived: input,
		Suggestion:    suggestion,
		Cause:         cause,
	}
}

func NewInternalError(message, input, suggestion string, cause error) *CLIError {
	return &CLIError{
		ExitCode:      ExitFailure,
		ErrorType:     "internal_error",
		Message:       message,
		InputReceived: input,
		Suggestion:    suggestion,
		Cause:         cause,
	}
}
