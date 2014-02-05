package gin

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	listener net.Listener
	proxy    *httputil.ReverseProxy
	builder  Builder
	runner   Runner
	to       *url.URL
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
	p.to = url

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
		p.runner.Run()
		if strings.ToLower(req.Header.Get("Upgrade")) == "websocket" {
			proxyWebsocket(res, req, p.to)
		} else {
			p.proxy.ServeHTTP(res, req)
		}
	}
}

func proxyWebsocket(w http.ResponseWriter, r *http.Request, host *url.URL) {
	d, err := net.Dial("tcp", host.Host)
	if err != nil {
		http.Error(w, "Error contacting backend server.", 500)
		fmt.Errorf("Error dialing websocket backend %s: %v", host, err)
		return
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Not a hijacker?", 500)
		return
	}
	nc, _, err := hj.Hijack()
	if err != nil {
		fmt.Errorf("Hijack error: %v", err)
		return
	}
	defer nc.Close()
	defer d.Close()

	err = r.Write(d)
	if err != nil {
		fmt.Errorf("Error copying request to target: %v", err)
		return
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}
	go cp(d, nc)
	go cp(nc, d)
	<-errc
}
