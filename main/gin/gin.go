package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/gin"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var startTime = time.Now()

var helpTemplate = "usage: {{.Name}} [-v|--version] [-h|--help] [--port=<port>] <url>\n"

func main() {
	// override the app help template
	cli.AppHelpTemplate = helpTemplate

	app := cli.NewApp()
	app.Name = "gin"
	app.Usage = "A development for go web apps"
	app.Action = MainAction
	app.Flags = []cli.Flag{
		cli.IntFlag{"port", 5678, "port for the proxy server"},
	}

	app.Run(os.Args)
}

func MainAction(c *cli.Context) {

  logger := log.New(os.Stdout, "[gin] ", 0)

	port := c.Int("port")
	appPort := strconv.Itoa(port + 1)

	os.Setenv("PORT", appPort)

	wd, err := os.Getwd()
	if err != nil {
    logger.Fatal(err)
	}

	builder := gin.NewBuilder(".")
	runner := gin.NewRunner(filepath.Join(wd, filepath.Base(wd)))
	runner.SetWriter(os.Stdout)
	proxy := gin.NewProxy(builder, runner)

	config := &gin.Config{
		Port:    port,
		ProxyTo: "http://localhost:" + appPort,
	}

	err = proxy.Run(config)
	if err != nil {
    logger.Fatal(err)
	}

  logger.Printf("listening on port %d\n", port)

	// build right now
	build(builder)

	// scan for changes
	scanChanges(func(path string) {
		build(builder)
	})
}

func build(builder gin.Builder) {
	err := builder.Build()
	if err != nil {
    fmt.Println(builder.Errors())
	}
	time.Sleep(100 * time.Millisecond)
}

type scanCallback func(path string)

func scanChanges(cb scanCallback) {
	for {
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			// TODO load ignore from config
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
