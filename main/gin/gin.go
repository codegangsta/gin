package main

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

var startTime = time.Now()

func main() {

	println("walking")

	for {
		scanChanges()
		time.Sleep(500 * time.Millisecond)
	}
}

func scanChanges() {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// TODO load ingnore from config
		if path == ".git" {
			return filepath.SkipDir
		}

		if info.ModTime().After(startTime) {
			println("Changes detected. Compiling...")
			startTime = time.Now()
			return errors.New("done")
		}

		return nil
	})
}
