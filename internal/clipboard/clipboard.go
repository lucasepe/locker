/*
Package clipboard provides macOS/Windows/  platform clipboard access

The most common operations are `Read` and `Write`. To use them:

	// write/read text format data of the clipboard, and
	// the byte buffer regarding the text are UTF8 encoded.
	clipboard.Write([]byte("text data"))
	clipboard.Read()

In addition, `clipboard.Write` returns a channel that can receive an
empty struct as a signal, which indicates the corresponding write call
to the clipboard is outdated, meaning the clipboard has been overwritten
by others and the previously written data is lost. For instance:

	changed := clipboard.Write([]byte("text data"))

	select {
	case <-changed:
		println(`"text data" is no longer available from clipboard.`)
	}

You can ignore the returning channel if you don't need this type of
notification. Furthermore, when you need more than just knowing whether
clipboard data is changed, use the watcher API:

	ch := clipboard.Watch(context.TODO())
	for data := range ch {
		// print out clipboard data whenever it is changed
		println(string(data))
	}
*/
package clipboard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	// activate only for running tests.
	debug             = false
	errUnavailable    = errors.New("clipboard unavailable")
	errUnsupported    = errors.New("unsupported format")
	errNotImplemented = errors.New("clipboard: cannot use when CGO_ENABLED=0 or os.arch is linux")
)

var (
	// Due to the limitation on operating systems (such as darwin),
	// concurrent read can even cause panic, use a global lock to
	// guarantee one read at a time.
	lock = sync.Mutex{}
)

// Read returns a chunk of bytes of the clipboard data if it presents
// in the desired format t presents. Otherwise, it returns nil.
func Read() []byte {
	lock.Lock()
	defer lock.Unlock()

	buf, err := read()
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "read clipboard err: %v\n", err)
		}
		return nil
	}
	return buf
}

// Write writes a given buffer to the clipboard.
// Write returned a receive-only channel can receive an empty struct
// as a signal, which indicates the clipboard has been overwritten from
// this write.
func Write(buf []byte) <-chan struct{} {
	lock.Lock()
	defer lock.Unlock()

	changed, err := write(buf)
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "write to clipboard err: %v\n", err)
		}
		return nil
	}
	return changed
}

// Watch returns a receive-only channel that received the clipboard data
// whenever any change of clipboard data in the desired format happens.
//
// The returned channel will be closed if the given context is canceled.
func Watch(ctx context.Context) <-chan []byte {
	return watch(ctx)
}
