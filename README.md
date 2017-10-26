gin [![wercker status](https://app.wercker.com/status/f413ccbd85cfc4a58a37f03dd7aaa87e "wercker status")](https://app.wercker.com/project/bykey/f413ccbd85cfc4a58a37f03dd7aaa87e)
========

`gin` is a simple command line utility for live-reloading Go web applications.
Just run `gin` in your app directory and your web app will be served with
`gin` as a proxy. `gin` will automatically recompile your code when it
detects a change. Your app will be restarted the next time it receives an
HTTP request.

`gin` adheres to the "silence is golden" principle, so it will only complain
if there was a compiler error or if you succesfully compile after an error.

## Installation

Assuming you have a working Go environment and `GOPATH/bin` is in your
`PATH`, `gin` is a breeze to install:

```shell
go get github.com/codegangsta/gin
```

Then verify that `gin` was installed correctly:

```shell
gin -h
```
## Basic usage
```shell
gin run main.go
```
Options
```
   --laddr value, -l value       listening address for the proxy server
   --port value, -p value        port for the proxy server (default: 3000)
   --appPort value, -a value     port for the Go web server (default: 3001)
   --bin value, -b value         name of generated binary file (default: "gin-bin")
   --path value, -t value        Path to watch files from (default: ".")
   --build value, -d value       Path to build files from (defaults to same value as --path)
   --excludeDir value, -x value  Relative directories to exclude
   --immediate, -i               run the server immediately after it's built
   --all                         reloads whenever any file changes, as opposed to reloading only on .go file change
   --godep, -g                   use godep when building
   --buildArgs value             Additional go build arguments
   --certFile value              TLS Certificate
   --keyFile value               TLS Certificate Key
   --logPrefix value             Setup custom log prefix
   --notifications               enable desktop notifications
   --help, -h                    show help
   --version, -v                 print the version
```

## Supporting Gin in Your Web app
`gin` assumes that your web app binds itself to the `PORT` environment
variable so it can properly proxy requests to your app. Web frameworks
like [Martini](http://github.com/codegangsta/martini) do this out of
the box.

## Using flags?
When you normally start your server with [flags](https://godoc.org/flag)
if you want to override any of them when running `gin` we suggest you
instead use [github.com/namsral/flag](https://github.com/namsral/flag)
as explained in [this post](http://stackoverflow.com/questions/24873883/organizing-environment-variables-golang/28160665#28160665)
