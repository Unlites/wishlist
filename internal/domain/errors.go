package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrForbidden = errors.New("forbidden")
var ErrUnauthorized = errors.New("unauthorized")
