package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/gin"
	"net/url"
	"os"
	"path/filepath"
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
	if len(c.Args()) != 1 {
		println("Error! Not enough arguments.\n")
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

  proxyToURL, err := url.Parse(c.Args()[0])
  if err != nil || len(proxyToURL.Host) == 0 {
		println("Error! Invalid URL")
		os.Exit(1)
  }
	port := c.Int("port")

	wd, err := os.Getwd()
	if err != nil {
		println(err)
		os.Exit(1)
	}

	builder := gin.NewBuilder(".")
	runner := gin.NewRunner(filepath.Join(wd, filepath.Base(wd)))
	runner.SetWriter(os.Stdout)
	proxy := gin.NewProxy(builder, runner)

	config := &gin.Config{
		Port:    port,
		ProxyTo: proxyToURL.String(),
	}

	println("gin server listening on port", port)
	err = proxy.Run(config)
	if err != nil {
		println(err)
		os.Exit(1)
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
