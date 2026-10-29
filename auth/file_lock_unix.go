//go:build !windows

package auth

import (
	"errors"
	"os"
	"syscall"
)

func tryLockExclusive(handle *os.File) (bool, error) {
	err := syscall.Flock(int(handle.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == nil {
		return true, nil
	}
	// EWOULDBLOCK and EAGAIN differ on some BSDs, so both count as contention, as does an interrupted call.
	if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EINTR) {
		return false, nil
	}
	return false, err
}

func unlockExclusive(handle *os.File) {
	_ = syscall.Flock(int(handle.Fd()), syscall.LOCK_UN)
}

// foldLockKey is how two spellings of one sidecar path are compared here: as written, unfolded.
func foldLockKey(absolutePath string) string { return absolutePath }
