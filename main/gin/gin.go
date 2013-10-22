package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/gin"
	"log"
	"os"
	"path/filepath"
	"time"
)

var startTime = time.Now()

func main() {
	app := cli.NewApp()
	app.Name = "gin"
	app.Usage = "A development server for go web apps"
	app.Action = MainAction

	app.Run(os.Args)
}

func MainAction(c *cli.Context) {
	println("Hello world")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	builder := gin.NewBuilder(".")
	runner := gin.NewRunner(filepath.Join(wd, filepath.Base(wd)))
	runner.SetWriter(os.Stdout)
	proxy := gin.NewProxy(builder, runner)

	config := &gin.Config{
		Port:    5678,
		ProxyTo: "http://localhost:3000",
	}

	err = proxy.Run(config)
	if err != nil {
		log.Fatalln(err)
	}

	// scan for changes
	scanChanges(func(path string) {
		println("building")
		err := builder.Build()
		if err != nil {
			println(builder.Errors())
		} else {
			println("Build successful")
		}
		time.Sleep(100 * time.Millisecond)
	})
}

type scanCallback func(path string)

func scanChanges(cb scanCallback) {
	for {
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			// TODO load ingnore from config
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
