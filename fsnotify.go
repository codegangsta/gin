//+build fsnotify
// Only build when the fsnotify tag is used

package main

import (
	"code.google.com/p/go.exp/fsnotify"
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
			walkPath(watchPath, cb)
		case err := <-watcher.Error:
			logger.Fatal(err)
		}
	}
}
