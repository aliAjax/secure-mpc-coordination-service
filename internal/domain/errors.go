package domain

import "errors"

var (
	ErrNotFound  = errors.New("mpc resource not found")
	ErrConflict  = errors.New("mpc state conflict")
	ErrInvalid   = errors.New("invalid mpc request")
	ErrLeaseLost = errors.New("round lease lost")
	ErrThreshold = errors.New("threshold not met")
	ErrReplay    = errors.New("replay rejected")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
