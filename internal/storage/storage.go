package storage

import "errors"

// здесь будет хранится общая информация

var (
	ErrURLNotFound = errors.New("URL not found")
	ErrVideoNotFound = errors.New("video not found")
	ErrURLExists = errors.New("URL already exists")
)