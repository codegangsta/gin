//+build !fsnotify
// The default is to build this scanner

package main

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

var (
	startTime = time.Now()
)

func scanChanges(watchPath string, cb scanCallback) {
	for {
		filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
			if path == ".git" {
				return filepath.SkipDir
			}

			if filepath.Ext(path) == ".go" && info.ModTime().After(startTime) {
				cb(path)
				startTime = time.Now()
				return errors.New("done")
			}

			return nil
		})

		time.Sleep(500 * time.Millisecond)
	}
}
