package util

import (
	"github.com/google/uuid"
)

// MustUUID - Create a UUID and throw a panic if an error occurs
func MustUUID() string {
	v, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return v
}

// NewUUID - Create a UUID
func NewUUID() (string, error) {
	v, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
