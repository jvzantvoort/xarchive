package main

import (
	"errors"
)

// Generic errors
var (
	// ErrEmptyList Empty list
	ErrEmptyList = errors.New("Empty list")
	// ErrNoTargets No targets in list
	ErrNoTargets = errors.New("No targets in list")
	// ErrNotFound Record not in list
	ErrNotFound = errors.New("Record not found")
	// ErrDuplicate Record is allready in list when it should not
	ErrDuplicate = errors.New("Record allready found")
)
