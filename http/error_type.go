package http

import (
	"slices"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
)

// ErrorType is an alias of the errors package's declaration, so the two spellings are one type.
type ErrorType = wcerrors.ErrorType

// The ErrorType values, re-exported from the errors package, which documents each one.
const (
	ErrorTypeNone             = wcerrors.ErrorTypeNone
	ErrorTypeRedirect         = wcerrors.ErrorTypeRedirect
	ErrorTypeValidation       = wcerrors.ErrorTypeValidation
	ErrorTypeAuthentication   = wcerrors.ErrorTypeAuthentication
	ErrorTypeAuthorization    = wcerrors.ErrorTypeAuthorization
	ErrorTypeNotFound         = wcerrors.ErrorTypeNotFound
	ErrorTypeMethodNotAllowed = wcerrors.ErrorTypeMethodNotAllowed
	ErrorTypeConflict         = wcerrors.ErrorTypeConflict
	ErrorTypeRateLimit        = wcerrors.ErrorTypeRateLimit
	ErrorTypeServer           = wcerrors.ErrorTypeServer
	ErrorTypeNetwork          = wcerrors.ErrorTypeNetwork
	ErrorTypeTimeout          = wcerrors.ErrorTypeTimeout
	ErrorTypeUnknown          = wcerrors.ErrorTypeUnknown
)

// In the order the refusal message names them.
var errorTypeMembers = []ErrorType{
	ErrorTypeNone, ErrorTypeRedirect, ErrorTypeValidation, ErrorTypeAuthentication,
	ErrorTypeAuthorization, ErrorTypeNotFound, ErrorTypeMethodNotAllowed, ErrorTypeConflict,
	ErrorTypeRateLimit, ErrorTypeServer, ErrorTypeNetwork, ErrorTypeTimeout,
	ErrorTypeUnknown,
}

func knownErrorType(errorType ErrorType) bool {
	return slices.Contains(errorTypeMembers, errorType)
}

// ClassifyErrorType maps an HTTP status to a classification.
func ClassifyErrorType(statusCode int) ErrorType {
	return classifyErrorType(statusCode, false)
}

// classifyErrorType classifies a 2xx whose body denies isSuccess as ErrorTypeUnknown; other statuses keep theirs.
func classifyErrorType(statusCode int, bodyDeniesSuccess bool) ErrorType {
	switch {
	case statusCode >= 200 && statusCode < 300:
		if bodyDeniesSuccess {
			return ErrorTypeUnknown
		}
		return ErrorTypeNone
	case statusCode >= 300 && statusCode < 400:
		return ErrorTypeRedirect
	case statusCode == 400 || statusCode == 422:
		return ErrorTypeValidation
	case statusCode == 401:
		return ErrorTypeAuthentication
	case statusCode == 403:
		return ErrorTypeAuthorization
	case statusCode == 404:
		return ErrorTypeNotFound
	case statusCode == 405:
		return ErrorTypeMethodNotAllowed
	case statusCode == 409:
		return ErrorTypeConflict
	case statusCode == 429:
		return ErrorTypeRateLimit
	case statusCode >= 500 && statusCode < 600:
		return ErrorTypeServer
	default:
		return ErrorTypeUnknown
	}
}
