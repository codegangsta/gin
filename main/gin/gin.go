package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"time"
)

func main() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Watch(pwd)
	if err != nil {
		log.Fatal(err)
	}

	func() {
		for {
			select {
			case <-watcher.Event:
				log.Println("Compiling, sleeping for 2s")
				time.Sleep(time.Second * 2)
				emptyWatcher(watcher)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	watcher.Close()
}

func emptyWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		case <-watcher.Event:
		default:
			return
		}
	}
}
