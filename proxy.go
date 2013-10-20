package gin

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	listener net.Listener
	proxy    *httputil.ReverseProxy
	builder  Builder
	runner   Runner
}

func NewProxy(builder Builder, runner Runner) *Proxy {
	return &Proxy{
		builder: builder,
		runner:  runner,
	}
}

func (p *Proxy) Run(config *Config) error {

	// create our reverse proxy
	url, err := url.Parse(config.ProxyTo)
	if err != nil {
		return err
	}
	p.proxy = httputil.NewSingleHostReverseProxy(url)

	p.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}

	go http.Serve(p.listener, http.HandlerFunc(p.defaultHandler))
	return nil
}

func (p *Proxy) Close() error {
	return p.listener.Close()
}

func (p *Proxy) defaultHandler(res http.ResponseWriter, req *http.Request) {
	errors := p.builder.Errors()
	if len(errors) > 0 {
		res.Write([]byte(errors))
	} else {
		p.proxy.ServeHTTP(res, req)
	}
}
