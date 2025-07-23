package errlist

import "errors"

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrBadRequest      = errors.New("bad request")
	ErrInternalServer  = errors.New("internal server error")

	// -----------------------------------------------------------

	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists          = errors.New("user already exists")
	ErrPasswordIsIncorrect = errors.New("password is incorrect")
	ErrUserIsNotVerified   = errors.New("user has not been verified")

	ErrNotAvailableSeats = errors.New("not available seats")
)
