package main

import (
	"github.com/codegangsta/cli"
	"github.com/howeyc/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	CompileErrors []byte
	Proxy         *httputil.ReverseProxy
	Dirty         = true
)

func main() {
	// set up logs
	log.SetFlags(0)
	log.SetPrefix("[passport] ")

	app := cli.NewApp()
	app.Name = "passport"
	app.Usage = "a Go development server"
	app.Action = DefaultAction

	app.Run(os.Args)
}

func DefaultAction(c *cli.Context) {
  url, err := url.Parse("http://localhost:3000")
  check(err)
  Proxy = httputil.NewSingleHostReverseProxy(url)

	go watch()
	go checkDirty()

	http.HandleFunc("/", MainHandler)
	err = http.ListenAndServe(":8080", nil)
	check(err)
}

func checkDirty() {
	var command *exec.Cmd
	for {
		if Dirty {
			log.Print("Restarting server...")
			build()
			if command != nil {
				command.Process.Kill()
			}
			command = run()
			Dirty = false
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	check(err)

	pwd, err := os.Getwd()
	check(err)

	err = watcher.Watch(pwd)
	check(err)

	defer watcher.Close()

	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsModify() && strings.HasSuffix(ev.Name, ".go") {
				Dirty = true
			}
		case err := <-watcher.Error:
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func MainHandler(res http.ResponseWriter, req *http.Request) {
	if len(CompileErrors) > 0 {
		res.Write(CompileErrors)
	} else {
		// TODO: Create proxy in main
		Proxy.ServeHTTP(res, req)
	}
}

func build() {
	command := exec.Command("go", "build")

	stderr, err := command.StderrPipe()
	check(err)

	err = command.Start()
	check(err)

	CompileErrors, err = ioutil.ReadAll(stderr)
	check(err)
}

func run() *exec.Cmd {
	wd, err := os.Getwd()
	check(err)

	command := exec.Command(filepath.Join(wd, filepath.Base(wd)))
	stdout, err := command.StdoutPipe()
	check(err)

	err = command.Start()
	check(err)

	go io.Copy(os.Stdout, stdout)

	return command
}

func check(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
