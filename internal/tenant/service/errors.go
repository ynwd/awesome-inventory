package service

import "errors"

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrTenantNotFound = errors.New("tenant not found")
	ErrOwnerNotFound  = errors.New("owner not found")
	ErrEmailExists    = errors.New("email already exists")
	ErrDatabaseExists = errors.New("database name already exists")
)
