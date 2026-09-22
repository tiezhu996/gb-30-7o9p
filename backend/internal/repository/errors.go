package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// Sentinel errors shared by all repositories.
var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

func isDuplicate(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "Duplicate entry"))
}

func translate(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}
