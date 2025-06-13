package task

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrFailedGetTask = errors.New("failed to get task")
	ErrUserNotFound  = errors.New("user ID not found in context")
)
