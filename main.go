package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/envy/lib"
	"github.com/codegangsta/gin/lib"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	logger     = log.New(os.Stdout, "[gin] ", 0)
	buildError error
)

type scanCallback func(path string)

func main() {
	app := cli.NewApp()
	app.Name = "gin"
	app.Usage = "A live reload utility for Go web applications."
	app.Action = MainAction
	app.Flags = []cli.Flag{
		cli.IntFlag{"port,p", 3000, "port for the proxy server"},
		cli.IntFlag{"appPort,a", 3001, "port for the Go web server"},
		cli.StringFlag{"bin,b", "gin-bin", "name of generated binary file"},
		cli.StringFlag{"path,t", ".", "Path to watch files from"},
	}
	app.Commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the gin proxy in the current working directory",
			Action:    MainAction,
		},
		{
			Name:      "env",
			ShortName: "e",
			Usage:     "Display environment variables set by the .env file",
			Action:    EnvAction,
		},
	}

	app.Run(os.Args)
}

func MainAction(c *cli.Context) {
	port := c.GlobalInt("port")
	appPort := strconv.Itoa(c.GlobalInt("appPort"))

	// Bootstrap the environment
	envy.Bootstrap()

	// Set the PORT env
	os.Setenv("PORT", appPort)

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}

	builder := gin.NewBuilder(".", c.GlobalString("bin"))
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
	scanChanges(c.GlobalString("path"), func(path string) {
		runner.Kill()
		build(builder, logger)
	})
}

func EnvAction(c *cli.Context) {
	// Bootstrap the environment
	env, err := envy.Bootstrap()
	if err != nil {
		logger.Fatalln(err)
	}

	for k, v := range env {
		fmt.Printf("%s: %s\n", k, v)
	}

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
