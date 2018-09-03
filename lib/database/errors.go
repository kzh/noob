package db

import (
	"errors"
)

var (
	ErrInvalidCredential   = errors.New("Invalid username or password.")
	ErrUnavailableUsername = errors.New("Username taken.")
	ErrInternalServer      = errors.New("Internal server error.")
	ErrNoSuchProblem       = errors.New("No such problem.")
)
