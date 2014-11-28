gin [![wercker status](https://app.wercker.com/status/f413ccbd85cfc4a58a37f03dd7aaa87e "wercker status")](https://app.wercker.com/project/bykey/f413ccbd85cfc4a58a37f03dd7aaa87e)
========

`gin` is a simple command line utility for live-reloading Go web applications. Just run `gin` in your app directory and your web app will be served with `gin` as a proxy. `gin` will automatically recompile your code when it detects a change. Your app will be restarted the next time it receives an HTTP request.

`gin` adheres to the "silence is golden" principle, so it will only complain if there was a compiler error or if you succesfully compile after an error.

## Installation

Assuming you have a working Go environment and `GOPATH/bin` is in your `PATH`, `gin` is a breeze to install:

~~~ shell
go get github.com/codegangsta/gin
~~~

Then verify that `gin` was installed correctly:

~~~ shell
gin -h
~~~

## Supporting Gin in Your Web app
`gin` assumes that your web app binds itself to the `PORT` environment variable so it can properly proxy requests to your app. Web frameworks like [Martini](http://github.com/codegangsta/martini) do this out of the box.
## Example app

```go 
package main

import (
  "fmt"
  "net/http"
  "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  http.HandleFunc("/", handler)

  port := os.Getenv("PORT")

  // add colon to fix error experienced on OSX
  // 2014/11/27 23:17:51 http: proxy error: dial tcp 127.0.0.1:3001: connection refused
  port = ":" + port

  fmt.Println(port)

  http.ListenAndServe(port, nil)
}
```