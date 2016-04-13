// +build !windows

package lock

//
// This code is inspired by https://github.com/boltdb/bolt
//
// The MIT License (MIT)
// Copyright (c) 2013 Ben Johnson
// https://github.com/boltdb/bolt/blob/master/LICENSE
//

import (
	"errors"
	"os"
	"syscall"
	"time"
)

// flock acquires an advisory lock on a file descriptor.
func Flock(f *os.File, exclusive bool, timeout time.Duration) error {
	var t time.Time
	for {
		// If we're beyond our timeout then return an error.
		// This can only occur after we've attempted a flock once.
		if t.IsZero() {
			t = time.Now()
		} else if timeout > 0 && time.Since(t) > timeout {
			return errors.New("timeout")
		}
		flag := syscall.LOCK_SH
		if exclusive {
			flag = syscall.LOCK_EX
		}

		// Otherwise attempt to obtain an exclusive lock.
		err := syscall.Flock(int(f.Fd()), flag|syscall.LOCK_NB)
		if err == nil {
			return nil
		} else if err != syscall.EWOULDBLOCK {
			return err
		}

		// Wait for a bit and try again.
		time.Sleep(50 * time.Millisecond)
	}
}

// funlock releases an advisory lock on a file descriptor.
func Funlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
