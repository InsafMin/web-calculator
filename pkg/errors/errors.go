package errors

import "errors"

var (
	ErrDivisionByZero       = errors.New("division by zero")
	ErrInvalidExpression    = errors.New("invalid expression")
	ErrOperatorNotSupported = errors.New("operator not supported")
	ErrUnacceptableSymbol   = errors.New("unacceptable symbol")
	ErrExtraOperator        = errors.New("extra operator")
	ErrExtraOpenBracket     = errors.New("extra open bracket")
	ErrExtraCloseBracket    = errors.New("extra close bracket")
	ErrNoTasksAvailable     = errors.New("no tasks available")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
