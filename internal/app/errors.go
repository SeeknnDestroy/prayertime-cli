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
	Details       *ErrorDetails
	Cause         error
}

type ErrorDetails struct {
	Candidates    []LocationCandidate `json:"candidates,omitempty"`
	ValidFields   []string            `json:"valid_fields,omitempty"`
	ValidTargets  []string            `json:"valid_targets,omitempty"`
	RequiredOneOf [][]string          `json:"required_one_of,omitempty"`
}

type LocationCandidate struct {
	DisplayName string  `json:"display_name"`
	CountryCode string  `json:"country_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
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

func (e *CLIError) WithDetails(details ErrorDetails) *CLIError {
	if e == nil {
		return nil
	}

	e.Details = &details
	return e
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
	return newNetworkCLIError("network_error", message, input, suggestion, cause)
}

func NewNetworkTimeoutError(message, input, suggestion string, cause error) *CLIError {
	return newNetworkCLIError("network_timeout", message, input, suggestion, cause)
}

func newNetworkCLIError(errorType, message, input, suggestion string, cause error) *CLIError {
	return &CLIError{
		ExitCode:      ExitNetwork,
		ErrorType:     errorType,
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
