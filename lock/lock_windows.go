// +build windows

package lock

import (
	"os"
	"time"
)

//
// This code is inspired by https://github.com/boltdb/bolt
//
// The MIT License (MIT)
// Copyright (c) 2013 Ben Johnson
// https://github.com/boltdb/bolt/blob/master/LICENSE
//

// flock acquires an advisory lock on a file descriptor.
func Flock(f *os.File, _ bool, _ time.Duration) error {
	return nil
}

// funlock releases an advisory lock on a file descriptor.
func Funlock(f *os.File) error {
	return nil
}
