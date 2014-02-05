package gin_test

import (
	"fmt"
	"github.com/codegangsta/gin/lib"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_NewProxy(t *testing.T) {
	builder := NewMockBuilder()
	runner := NewMockRunner()
	proxy := gin.NewProxy(builder, runner)

	expect(t, proxy != nil, true)
}

func Test_Proxy_Run(t *testing.T) {
	builder := NewMockBuilder()
	runner := NewMockRunner()
	proxy := gin.NewProxy(builder, runner)

	config := &gin.Config{}

	proxy.Run(config)
	defer proxy.Close()
}

func Test_Proxying(t *testing.T) {
	builder := NewMockBuilder()
	runner := NewMockRunner()
	proxy := gin.NewProxy(builder, runner)

	// create a test server and see if we can proxy a request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	}))
	defer ts.Close()

	config := &gin.Config{
		Port:    5678,
		ProxyTo: ts.URL,
	}

	err := proxy.Run(config)
	defer proxy.Close()
	expect(t, err, nil)

	res, err := http.Get("http://localhost:5678")
	expect(t, err, nil)
	expect(t, res == nil, false)
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	expect(t, fmt.Sprintf("%s", greeting), "Hello world\n")
	expect(t, runner.DidRun, true)
}

func Test_Proxying_Websocket(t *testing.T) {
	builder := NewMockBuilder()
	runner := NewMockRunner()
	proxy := gin.NewProxy(builder, runner)

	// create a test server and see if we can proxy a websocket request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	}))
	defer ts.Close()

	config := &gin.Config{
		Port:    5678,
		ProxyTo: ts.URL,
	}

	err := proxy.Run(config)
	defer proxy.Close()
	expect(t, err, nil)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:5678", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "Websocket")
	res, _ := client.Do(req)
	expect(t, err, nil)
	expect(t, res == nil, false)
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	expect(t, fmt.Sprintf("%s", greeting), "Hello world\n")
	expect(t, runner.DidRun, true)
}

func Test_Proxying_Build_Errors(t *testing.T) {
	builder := NewMockBuilder()
	builder.MockErrors = "Foo bar here are some errors"
	runner := NewMockRunner()
	proxy := gin.NewProxy(builder, runner)

	config := &gin.Config{
		Port:    5679,
		ProxyTo: "http://localhost:3000",
	}

	err := proxy.Run(config)
	defer proxy.Close()
	expect(t, err, nil)

	res, err := http.Get("http://localhost:5679")
	expect(t, err, nil)
	expect(t, res == nil, false)
	errors, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	expect(t, fmt.Sprintf("%s", errors), "Foo bar here are some errors")
}
