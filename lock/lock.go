package lock

//
// This code is inspired by https://github.com/boltdb/bolt
//
// The MIT License (MIT)
// Copyright (c) 2013 Ben Johnson
// https://github.com/boltdb/bolt/blob/master/LICENSE
//

import (
	"os"
	"time"
)

type LockFile struct {
	file *os.File
}

func Lock(filePath string, exclusive bool, timeout time.Duration) (*LockFile, error) {
	flag := os.O_RDWR
	if !exclusive {
		flag = os.O_RDONLY
	}

	// Open lock file.
	fh, err := os.OpenFile(filePath, flag|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	if err := Flock(fh, exclusive, timeout); err != nil {
		return nil, err
	}

	return &LockFile{
		file: fh,
	}, nil
}

func (l *LockFile) Path() string {
	return l.file.Name()
}

func (l *LockFile) UnLock() error {
	return Funlock(l.file)
}

func (l *LockFile) Remove() error {
	return os.Remove(l.file.Name())
}
