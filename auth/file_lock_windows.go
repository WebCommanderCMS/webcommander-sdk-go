//go:build windows

package auth

import (
	"errors"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// The two calls come from kernel32 directly: syscall exports neither, and go.mod stays dependency-free.
var (
	kernel32     = syscall.NewLazyDLL("kernel32.dll")
	procLockFile = kernel32.NewProc("LockFileEx")
	procUnlock   = kernel32.NewProc("UnlockFileEx")
)

const (
	lockfileFailImmediately = 0x00000001
	lockfileExclusiveLock   = 0x00000002
	// ERROR_LOCK_VIOLATION: somebody else holds an overlapping range right now.
	errorLockViolation = syscall.Errno(33)
)

func tryLockExclusive(handle *os.File) (bool, error) {
	overlapped := syscall.Overlapped{}
	result, _, err := procLockFile.Call(
		handle.Fd(),
		uintptr(lockfileExclusiveLock|lockfileFailImmediately),
		0,
		1,                                    // one byte, low word
		0,                                    // one byte, high word
		uintptr(unsafe.Pointer(&overlapped)), //#nosec G103 -- an LPOVERLAPPED argument, a pointer-free local, converted inline for a synchronous call
	)
	if result != 0 {
		return true, nil
	}
	if errors.Is(err, errorLockViolation) {
		return false, nil
	}
	return false, err
}

func unlockExclusive(handle *os.File) {
	overlapped := syscall.Overlapped{}
	_, _, _ = procUnlock.Call(handle.Fd(), 0, 1, 0,
		uintptr(unsafe.Pointer(&overlapped))) //#nosec G103 -- an LPOVERLAPPED argument, a pointer-free local, converted inline for a synchronous call
}

// foldLockKey lower-cases a sidecar path, as Windows compares paths case-insensitively.
func foldLockKey(absolutePath string) string { return strings.ToLower(absolutePath) }
