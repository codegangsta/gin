//+build !fsnotify
// The default is to build this scanner

package main

import (
	"time"
)

func scanChanges(watchPath string, cb scanCallback) {
	for {
		walkPath(watchPath, cb)
		time.Sleep(500 * time.Millisecond)
	}
}
