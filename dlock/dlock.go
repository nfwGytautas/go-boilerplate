// Package dlock provides a mechanism for distributed locking this is used to
// synchronize services with multiple instances, do not this should be used
// sparingly, as there aren't that many cases that couldn't be solved with
// smart design
package dlock

import (
	"context"
	"errors"
)

var (
	ErrLocked          = errors.New("lock already acquired")
	ErrNonOwnerRelease = errors.New("trying to release lock with acquiring it first")
)

// DLock interface definining a locking mechanism interface
type DLock interface {
	// Acquire locks the section
	Acquire(context.Context) error

	// Release releases the section lock
	Release(context.Context) error
}
