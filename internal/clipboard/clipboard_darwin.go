//go:build darwin && !ios
// +build darwin,!ios

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

unsigned int clipboard_read_string(void **out);
unsigned int clipboard_read_image(void **out);
int clipboard_write_string(const void *bytes, NSInteger n);
int clipboard_write_image(const void *bytes, NSInteger n);
NSInteger clipboard_change_count();
*/
import "C"
import (
	"context"
	"time"
	"unsafe"
)

func read() (buf []byte, err error) {
	var (
		data unsafe.Pointer
		n    C.uint
	)
	n = C.clipboard_read_string(&data)
	if data == nil {
		return nil, errUnavailable
	}
	defer C.free(unsafe.Pointer(data))
	if n == 0 {
		return nil, nil
	}
	return C.GoBytes(data, C.int(n)), nil
}

// write writes the given data to clipboard and
// returns true if success or false if failed.
func write(buf []byte) (bool, error) {
	var ok C.int

	if len(buf) == 0 {
		ok = C.clipboard_write_string(unsafe.Pointer(nil), 0)
	} else {
		ok = C.clipboard_write_string(unsafe.Pointer(&buf[0]),
			C.NSInteger(len(buf)))
	}
	if ok != 0 {
		return false, errUnavailable
	}

	return true, nil
}

func watch(ctx context.Context) <-chan []byte {
	recv := make(chan []byte, 1)
	// not sure if we are too slow or the user too fast :)
	ti := time.NewTicker(time.Second)
	lastCount := C.long(C.clipboard_change_count())
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(recv)
				return
			case <-ti.C:
				this := C.long(C.clipboard_change_count())
				if lastCount != this {
					b := Read()
					if b == nil {
						continue
					}
					recv <- b
					lastCount = this
				}
			}
		}
	}()
	return recv
}
