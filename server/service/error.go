package service

import "fmt"

type ErrorCode uint

const (
	ErrGeneric ErrorCode = iota
	ErrInvalidMsgFormat
	ErrInvalidPayloadFormat
	ErrInvalidAction
	ErrChallengeAlreadyProvided
	ErrInvalidSolution
	ErrNoActiveChallenge
)

var codeStr = []string{
	ErrGeneric:                  "ErrGeneric",
	ErrInvalidMsgFormat:         "ErrInvalidMsgFormat",
	ErrInvalidPayloadFormat:     "ErrInvalidPayloadFormat",
	ErrInvalidAction:            "ErrInvalidAction",
	ErrChallengeAlreadyProvided: "ErrChallengeAlreadyProvided",
	ErrInvalidSolution:          "ErrInvalidSolution",
	ErrNoActiveChallenge:        "ErrNoActiveChallenge",
}

var msg = map[ErrorCode]string{
	ErrGeneric:                  "something went wrong",
	ErrInvalidMsgFormat:         "invalid message format",
	ErrInvalidPayloadFormat:     "invalid payload format",
	ErrInvalidAction:            "invalid action",
	ErrChallengeAlreadyProvided: "challenge has been already provided",
	ErrInvalidSolution:          "solution is invalid",
	ErrNoActiveChallenge:        "no active challenge for token",
}

func (ec ErrorCode) String() string {
	return codeStr[ec]
}

type Error struct {
	code    ErrorCode
	message string
}

// NewError creates Error with code and original error. If original error
// is nil NewError returns nil
func NewError(code ErrorCode, originalError error) *Error {
	if originalError == nil {
		return nil
	}

	message := fmt.Errorf("%s: %w", msg[code], originalError).Error()

	return &Error{
		code:    code,
		message: message,
	}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Code() ErrorCode {
	return e.code
}
