//+build fsnotify
// Only build when the fsnotify tag is used

package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"errors"
	"os"
	"path/filepath"
	"time"
)

func scanChanges(watchPath string, cb scanCallback) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}

	err = watcher.Watch(watchPath)

	if err != nil {
		logger.Fatal(err)
	}

	for {
		select {
		case <-watcher.Event:
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
		case err := <-watcher.Error:
			logger.Fatal(err)
		}
	}
}
