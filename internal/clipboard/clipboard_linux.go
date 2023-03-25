//go:build linux && !cgo
// +build linux,!cgo

package clipboard

import "context"

func read() (buf []byte, err error) {
	return nil, errNotImplemented
}

func write() (bool, error) {
	return false, errNotImplemented
}

func watch(ctx context.Context) <-chan []byte {
	return nil, errNotImplemented
}
