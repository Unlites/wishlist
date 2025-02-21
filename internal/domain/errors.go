package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrAlreadyExists    = errors.New("already exists")
	ErrAlreadyProcessed = errors.New("already processed")
	ErrForbidden        = errors.New("forbidden")
)
