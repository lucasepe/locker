//go:build linux && !cgo
// +build linux,!cgo

package clipboard

import "context"

func read() (buf []byte, err error) {
	return nil, errNotImplemented
}

func write() (<-chan struct{}, error) {
	return nil, errNotImplemented
}

func watch(ctx context.Context) <-chan []byte {
	return nil, errNotImplemented
}
