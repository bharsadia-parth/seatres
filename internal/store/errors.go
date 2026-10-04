package store

import "errors"

var (
	ErrNotFound    = errors.New("not_found")
	ErrForbidden   = errors.New("forbidden")
	ErrSeatTaken   = errors.New("seat_taken")
	ErrLimit       = errors.New("per_user_limit")
	ErrKeyMismatch = errors.New("idempotency_mismatch")
)
