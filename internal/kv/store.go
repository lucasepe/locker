package kv

import (
	"errors"
)

const (
	EnvSecret = "LOCKER_SECRET"
)

var (
	ErrEmptyNamespace      = errors.New("namespace cannot be empty")
	ErrEmptyKey            = errors.New("key cannot be empty")
	ErrUnsetMasterPassword = errors.New("master password cannot be empty")
	ErrNamespaceNotFound   = errors.New("namespace not found")
)

// Store is an abstraction for different key-value store implementations.
// A store must be able to store, retrieve and delete key-value pairs,
// in the specified namespace.
type Store interface {
	// PutOne stores the given value for the given key in the specified namespace.
	PutOne(namespace string, key, value string) error
	// GetOne retrieves the value for the given key in the specified namespace.
	GetOne(namespace string, key string) (val string, err error)
	// GetAll retrieves all the values in a given namespace
	GetAll(namespace string) (map[string]string, error)
	// DeleteOne deletes the stored value for the given key in the specified namespace.
	DeleteOne(namespace string, key string) error
	// DeleteAll deletes all the values in a namespace.
	DeleteAll(namespace string) error
	// Namespaces returns all namespace names.
	Namespaces() (names []string, err error)
	// Keys returns all keys in a namespace.
	Keys(namespace string) (items []string, err error)
	// Close must be called when the work with the key-value store is done.
	Close() error
}
