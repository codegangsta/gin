//+build fsnotify
// Only build when the fsnotify tag is used

package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"os"
	"path/filepath"
	"time"
)

// This is to prevent unneeded rebuilds
var (
	watched = make(map[string]bool)
)

func addPaths(watcher *fsnotify.Watcher, path string) error {
	if watched[path] {
		return nil
	}

	return filepath.Walk(path, func(wpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if wpath == ".git" {
			return filepath.SkipDir
		}

		if watched[wpath] {
			return nil
		}

		if info.IsDir() {
			watched[wpath] = true
			err = watcher.Watch(wpath)
			if err != nil {
				return err
			}
			if path != wpath {
				return addPaths(watcher, wpath)
			}
		} else {
			if filepath.Ext(wpath) == ".go" {
				watched[wpath] = true
			}
		}

		return nil
	})
}

func scanChanges(watchPath string, cb scanCallback) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}

	err = addPaths(watcher, watchPath)
	if err != nil {
		logger.Fatal(err)
	}

	for {
		select {
		case ev := <-watcher.Event:
			path := ev.Name
			if watched[path] {
				if ev.IsDelete() {
					err = watcher.RemoveWatch(path)
					if err != nil {
						logger.Fatal(err)
					}
					watched[path] = false
					continue
				}

				if filepath.Ext(path) == ".go" {
					cb(path)
					// Drain events that might have been triggered since this build
					drain := true
					for drain {
						select {
						case <-watcher.Event:
						case <-time.After(250 * time.Millisecond):
							drain = false
						}
					}
				}
			} else {
				if ev.IsCreate() {
					// We don't really care /too/ much if this fails...
					addPaths(watcher, path)
				}
			}

		case err := <-watcher.Error:
			logger.Fatal(err)
		}
	}
}
