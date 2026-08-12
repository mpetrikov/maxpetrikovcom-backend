package lab

import "errors"

var (
	ErrNotFound          = errors.New("lab not found")
	ErrSlugAlreadyExists = errors.New("lab slug already exists")
)
