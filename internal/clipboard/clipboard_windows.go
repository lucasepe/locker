// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows
// +build windows

package clipboard

// Interacting with Clipboard on Windows:
// https://learn.microsoft.com/en-us/windows/win32/dataxchg/using-the-clipboard

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"reflect"
	"runtime"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

// readText reads the clipboard and returns the text data if presents.
// The caller is responsible for opening/closing the clipboard before
// calling this function.
func readText() (buf []byte, err error) {
	hMem, _, err := getClipboardData.Call(cFmtUnicodeText)
	if hMem == 0 {
		return nil, err
	}
	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return nil, err
	}
	defer gUnlock.Call(hMem)

	// Find NUL terminator
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint16)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) +
			unsafe.Sizeof(*((*uint16)(unsafe.Pointer(p)))))
	}

	var s []uint16
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Data = p
	h.Len = n
	h.Cap = n
	return []byte(string(utf16.Decode(s))), nil
}

// writeText writes given data to the clipboard. It is the caller's
// responsibility for opening/closing the clipboard before calling
// this function.
func writeText(buf []byte) error {
	r, _, err := emptyClipboard.Call()
	if r == 0 {
		return fmt.Errorf("failed to clear clipboard: %w", err)
	}

	// empty text, we are done here.
	if len(buf) == 0 {
		return nil
	}

	s, err := syscall.UTF16FromString(string(buf))
	if err != nil {
		return fmt.Errorf("failed to convert given string: %w", err)
	}

	hMem, _, err := gAlloc.Call(gmemMoveable, uintptr(len(s)*int(unsafe.Sizeof(s[0]))))
	if hMem == 0 {
		return fmt.Errorf("failed to alloc global memory: %w", err)
	}

	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return fmt.Errorf("failed to lock global memory: %w", err)
	}
	defer gUnlock.Call(hMem)

	// no return value
	memMove.Call(p, uintptr(unsafe.Pointer(&s[0])),
		uintptr(len(s)*int(unsafe.Sizeof(s[0]))))

	v, _, err := setClipboardData.Call(cFmtUnicodeText, hMem)
	if v == 0 {
		gFree.Call(hMem)
		return fmt.Errorf("failed to set text to clipboard: %w", err)
	}

	return nil
}

func read() (buf []byte, err error) {
	// On Windows, OpenClipboard and CloseClipboard must be executed on
	// the same thread. Thus, lock the OS thread for further execution.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// check if clipboard is avaliable for the requested format
	r, _, err := isClipboardFormatAvailable.Call(cFmtUnicodeText)
	if r == 0 {
		return nil, errUnavailable
	}

	// try again until open clipboard successed
	for {
		r, _, _ = openClipboard.Call()
		if r == 0 {
			continue
		}
		break
	}
	defer closeClipboard.Call()

	switch format {
	case cFmtDIBV5:
		return nil, errUnsupported
	case cFmtUnicodeText:
		fallthrough
	default:
		return readText()
	}
}

// write writes the given data to clipboard and
// returns true if success or false if failed.
func write(buf []byte) (bool, error) {
	errch := make(chan error)
	changed := make(chan struct{}, 1)
	go func() {
		// make sure GetClipboardSequenceNumber happens with
		// OpenClipboard on the same thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		for {
			r, _, _ := openClipboard.Call(0)
			if r == 0 {
				continue
			}
			break
		}

		// param = cFmtUnicodeText
		err := writeText(buf)
		if err != nil {
			errch <- err
			closeClipboard.Call()
			return
		}

		// Close the clipboard otherwise other applications cannot
		// paste the data.
		closeClipboard.Call()

		cnt, _, _ := getClipboardSequenceNumber.Call()
		errch <- nil
		for {
			time.Sleep(time.Second)
			cur, _, _ := getClipboardSequenceNumber.Call()
			if cur != cnt {
				changed <- struct{}{}
				close(changed)
				return
			}
		}
	}()
	err := <-errch
	if err != nil {
		return false, err
	}
	_, ok := <-changed
	return ok, nil
}

func watch(ctx context.Context) <-chan []byte {
	recv := make(chan []byte, 1)
	ready := make(chan struct{})
	go func() {
		// not sure if we are too slow or the user too fast :)
		ti := time.NewTicker(time.Second)
		cnt, _, _ := getClipboardSequenceNumber.Call()
		ready <- struct{}{}
		for {
			select {
			case <-ctx.Done():
				close(recv)
				return
			case <-ti.C:
				cur, _, _ := getClipboardSequenceNumber.Call()
				if cnt != cur {
					b := Read(t)
					if b == nil {
						continue
					}
					recv <- b
					cnt = cur
				}
			}
		}
	}()
	<-ready
	return recv
}

const (
	cFmtBitmap      = 2 // Win+PrintScreen
	cFmtUnicodeText = 13
	cFmtDIBV5       = 17
	// Screenshot taken from special shortcut is in different format (why??), see:
	// https://jpsoft.com/forums/threads/detecting-clipboard-format.5225/
	cFmtDataObject = 49161 // Shift+Win+s, returned from enumClipboardFormats
	gmemMoveable   = 0x0002
)

// Calling a Windows DLL, see:
// https://github.com/golang/go/wiki/WindowsDLLs
var (
	user32 = syscall.MustLoadDLL("user32")
	// Opens the clipboard for examination and prevents other
	// applications from modifying the clipboard content.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-openclipboard
	openClipboard = user32.MustFindProc("OpenClipboard")
	// Closes the clipboard.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-closeclipboard
	closeClipboard = user32.MustFindProc("CloseClipboard")
	// Empties the clipboard and frees handles to data in the clipboard.
	// The function then assigns ownership of the clipboard to the
	// window that currently has the clipboard open.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-emptyclipboard
	emptyClipboard = user32.MustFindProc("EmptyClipboard")
	// Retrieves data from the clipboard in a specified format.
	// The clipboard must have been opened previously.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getclipboarddata
	getClipboardData = user32.MustFindProc("GetClipboardData")
	// Places data on the clipboard in a specified clipboard format.
	// The window must be the current clipboard owner, and the
	// application must have called the OpenClipboard function. (When
	// responding to the WM_RENDERFORMAT message, the clipboard owner
	// must not call OpenClipboard before calling SetClipboardData.)
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setclipboarddata
	setClipboardData = user32.MustFindProc("SetClipboardData")
	// Determines whether the clipboard contains data in the specified format.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-isclipboardformatavailable
	isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable")
	// Clipboard data formats are stored in an ordered list. To perform
	// an enumeration of clipboard data formats, you make a series of
	// calls to the EnumClipboardFormats function. For each call, the
	// format parameter specifies an available clipboard format, and the
	// function returns the next available clipboard format.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-isclipboardformatavailable
	enumClipboardFormats = user32.MustFindProc("EnumClipboardFormats")
	// Retrieves the clipboard sequence number for the current window station.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getclipboardsequencenumber
	getClipboardSequenceNumber = user32.MustFindProc("GetClipboardSequenceNumber")
	// Registers a new clipboard format. This format can then be used as
	// a valid clipboard format.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-registerclipboardformata
	registerClipboardFormatA = user32.MustFindProc("RegisterClipboardFormatA")

	kernel32 = syscall.NewLazyDLL("kernel32")

	// Locks a global memory object and returns a pointer to the first
	// byte of the object's memory block.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globallock
	gLock = kernel32.NewProc("GlobalLock")
	// Decrements the lock count associated with a memory object that was
	// allocated with GMEM_MOVEABLE. This function has no effect on memory
	// objects allocated with GMEM_FIXED.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalunlock
	gUnlock = kernel32.NewProc("GlobalUnlock")
	// Allocates the specified number of bytes from the heap.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalalloc
	gAlloc = kernel32.NewProc("GlobalAlloc")
	// Frees the specified global memory object and invalidates its handle.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalfree
	gFree   = kernel32.NewProc("GlobalFree")
	memMove = kernel32.NewProc("RtlMoveMemory")
)
