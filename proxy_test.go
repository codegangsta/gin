package gin_test

import (
	"fmt"
	"github.com/codegangsta/gin"
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
}
