package services

import "errors"

var (
	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")
	
	// ErrUnauthorized is returned when user is not authorized
	ErrUnauthorized = errors.New("unauthorized")
	
	// ErrForbidden is returned when user doesn't have permission
	ErrForbidden = errors.New("forbidden")
) 