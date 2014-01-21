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

var (
	startTime    = time.Now()
	helpTemplate = "gin - live reload for martini\nusage: {{.Name}} [-v|--version] [-h|--help] [(-p|--port)=<port>] [(-a|--appPort)=<appPort>] [(-b|--bin)=<binary>]\n"
	logger       = log.New(os.Stdout, "[gin] ", 0)
	buildError   error
)

func main() {
	// override the app help template
	cli.AppHelpTemplate = helpTemplate

	app := cli.NewApp()
	app.Name = "gin"
	app.Usage = "A development server for martini"
	app.Action = MainAction
	app.Flags = []cli.Flag{
		cli.IntFlag{"port,p", 3000, "port for the proxy server"},
		cli.IntFlag{"appPort,a", 3001, "port for the martini web server"},
		cli.StringFlag{"bin,b", "gin-bin", "name of generated binary file"},
	}

	app.Run(os.Args)
}

func MainAction(c *cli.Context) {
	port := c.Int("port")
	appPort := strconv.Itoa(c.Int("appPort"))
	os.Setenv("PORT", appPort)

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}

	builder := gin.NewBuilder(".", c.String("bin"))
	runner := gin.NewRunner(filepath.Join(wd, builder.Binary()))
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
	build(builder, logger)

	// scan for changes
	scanChanges(func(path string) {
		runner.Kill()
		build(builder, logger)
	})
}

func build(builder gin.Builder, logger *log.Logger) {
	err := builder.Build()
	if err != nil {
		buildError = err
		logger.Println("ERROR! Build failed.")
		fmt.Println(builder.Errors())
	} else {
		// print success only if there were errors before
		if buildError != nil {
			logger.Println("Build Successful")
		}
		buildError = nil
	}

	time.Sleep(100 * time.Millisecond)
}

type scanCallback func(path string)

func scanChanges(cb scanCallback) {
	for {
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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
